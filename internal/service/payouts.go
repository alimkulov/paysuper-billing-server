package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/now"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/paysuper/paysuper-proto/go/postmarkpb"
	"github.com/paysuper/paysuper-proto/go/reporterpb"
	tools "github.com/paysuper/paysuper-tools/number"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"mime"
	"path/filepath"
	"sort"
	"time"
)

const (
	payoutChangeSourceMerchant = "merchant"
	payoutChangeSourceAdmin    = "admin"

	payoutArrivalInDays = 5
)

var (
	errorPayoutSourcesNotFound         = newBillingServerErrorMsg("po000001", "no source documents found for payout")
	errorPayoutSourcesPending          = newBillingServerErrorMsg("po000002", "you have at least one royalty report waiting for acceptance")
	errorPayoutSourcesDispute          = newBillingServerErrorMsg("po000003", "you have at least one unclosed dispute in your royalty reports")
	errorPayoutNotFound                = newBillingServerErrorMsg("po000004", "payout document not found")
	errorPayoutAmountInvalid           = newBillingServerErrorMsg("po000005", "payout amount is invalid")
	errorPayoutUpdateBalance           = newBillingServerErrorMsg("po000008", "balance update failed")
	errorPayoutBalanceError            = newBillingServerErrorMsg("po000009", "getting balance failed")
	errorPayoutNotEnoughBalance        = newBillingServerErrorMsg("po000010", "not enough balance for payout")
	errorPayoutUpdateRoyaltyReports    = newBillingServerErrorMsg("po000012", "royalty reports update failed")
	errorPayoutStatusChangeIsForbidden = newBillingServerErrorMsg("po000014", "status change is forbidden")
	errorPayoutManualPayoutsDisabled   = newBillingServerErrorMsg("po000015", "manual payouts disabled")
	errorPayoutAutoPayoutsDisabled     = newBillingServerErrorMsg("po000016", "auto payouts disabled")
	errorPayoutAutoPayoutsWithErrors   = newBillingServerErrorMsg("po000017", "auto payouts creation finished with errors")

	statusForUpdateBalance = map[string]bool{
		pkg.PayoutDocumentStatusPending: true,
		pkg.PayoutDocumentStatusPaid:    true,
	}

	statusForBecomePaid = map[string]bool{
		pkg.PayoutDocumentStatusPaid: true,
	}

	statusForBecomeFailed = map[string]bool{
		pkg.PayoutDocumentStatusFailed:   true,
		pkg.PayoutDocumentStatusCanceled: true,
	}
)

func (s *Service) CreatePayoutDocument(
	ctx context.Context,
	req *billingpb.CreatePayoutDocumentRequest,
	res *billingpb.CreatePayoutDocumentResponse,
) error {

	merchant, err := s.merchantRepository.GetById(ctx, req.MerchantId)
	if err != nil {
		return merchantErrorNotFound
	}

	if merchant.ManualPayoutsEnabled == req.IsAutoGeneration {
		res.Status = billingpb.ResponseStatusBadData
		if req.IsAutoGeneration {
			res.Message = errorPayoutAutoPayoutsDisabled
		} else {
			res.Message = errorPayoutManualPayoutsDisabled
		}
		return nil
	}

	return s.createPayoutDocument(ctx, merchant, req, res)
}

func (s *Service) createPayoutDocument(
	ctx context.Context,
	merchant *billingpb.Merchant,
	req *billingpb.CreatePayoutDocumentRequest,
	res *billingpb.CreatePayoutDocumentResponse,
) error {

	arrivalDate, err := ptypes.TimestampProto(now.EndOfDay().Add(time.Hour * 24 * payoutArrivalInDays))
	if err != nil {
		return err
	}

	pd := &billingpb.PayoutDocument{
		Id:                      primitive.NewObjectID().Hex(),
		Status:                  pkg.PayoutDocumentStatusPending,
		SourceId:                []string{},
		Description:             req.Description,
		CreatedAt:               ptypes.TimestampNow(),
		UpdatedAt:               ptypes.TimestampNow(),
		ArrivalDate:             arrivalDate,
		MerchantId:              merchant.Id,
		Destination:             merchant.Banking,
		Company:                 merchant.Company,
		MerchantAgreementNumber: merchant.AgreementNumber,
		OperatingCompanyId:      merchant.OperatingCompanyId,
	}

	reports, err := s.getPayoutDocumentSources(ctx, merchant)

	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	pd.Currency = reports[0].Currency

	var times []time.Time

	for _, r := range reports {
		pd.TotalFees += r.Totals.PayoutAmount - r.Totals.CorrectionAmount
		pd.Balance += r.Totals.PayoutAmount - r.Totals.CorrectionAmount - r.Totals.RollingReserveAmount
		pd.TotalTransactions += r.Totals.TransactionsCount
		pd.SourceId = append(pd.SourceId, r.Id)

		from, err := ptypes.Timestamp(r.PeriodFrom)

		if err != nil {
			zap.L().Error(
				"Time conversion error",
				zap.Error(err),
			)
			return err
		}

		to, err := ptypes.Timestamp(r.PeriodTo)
		if err != nil {
			zap.L().Error(
				"Payout source time conversion error",
				zap.Error(err),
			)
			return err
		}
		times = append(times, from, to)
	}

	if pd.Balance <= 0 {
		res.Status = billingpb.ResponseStatusBadData
		res.Message = errorPayoutAmountInvalid
		return nil
	}

	balance, err := s.getMerchantBalance(ctx, merchant.Id)
	if err != nil {
		res.Status = billingpb.ResponseStatusSystemError
		res.Message = errorPayoutBalanceError
		return nil
	}

	balanceNet := tools.ToPrecise(balance.Debit - balance.Credit)

	if pd.Balance > balanceNet {
		res.Status = billingpb.ResponseStatusBadData
		res.Message = errorPayoutNotEnoughBalance
		return nil
	}

	if pd.Balance < merchant.MinPayoutAmount {
		pd.Status = pkg.PayoutDocumentStatusSkip
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	from := times[0]
	to := times[len(times)-1]

	pd.PeriodFrom, err = ptypes.TimestampProto(from)
	if err != nil {
		zap.L().Error(
			"Payout PeriodFrom time conversion error",
			zap.Error(err),
		)
		return err
	}
	pd.PeriodTo, err = ptypes.TimestampProto(to)
	if err != nil {
		zap.L().Error(
			"Payout PeriodTo time conversion error",
			zap.Error(err),
		)
		return err
	}

	err = s.payoutRepository.Insert(ctx, pd, req.Ip, payoutChangeSourceMerchant)

	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusSystemError
			res.Message = e
			return nil
		}
		return err
	}

	err = s.royaltyReportSetPayoutDocumentId(ctx, pd.SourceId, pd.Id, req.Ip, req.Initiator)

	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusSystemError
			res.Message = e
			return nil
		}
		return err
	}

	err = s.renderPayoutDocument(ctx, pd, merchant)
	if err != nil {
		return err
	}

	res.Items = append(res.Items, pd)

	res.Status = billingpb.ResponseStatusOk

	return nil
}

func (s *Service) GetPayoutDocument(
	ctx context.Context,
	req *billingpb.GetPayoutDocumentRequest,
	res *billingpb.PayoutDocumentResponse,
) (err error) {
	res.Item, err = s.payoutRepository.GetByIdMerchantId(ctx, req.PayoutDocumentId, req.MerchantId)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPayoutNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Status = billingpb.ResponseStatusOk

	return nil
}

func (s *Service) GetPayoutDocumentRoyaltyReports(
	ctx context.Context,
	req *billingpb.GetPayoutDocumentRequest,
	res *billingpb.ListRoyaltyReportsResponse,
) error {
	pd, err := s.payoutRepository.GetByIdMerchantId(ctx, req.PayoutDocumentId, req.MerchantId)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPayoutNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Data = &billingpb.RoyaltyReportsPaginate{}
	res.Data.Items, err = s.royaltyReportRepository.GetByPayoutId(ctx, pd.Id)

	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Data.Count = int64(len(res.Data.Items))
	res.Status = billingpb.ResponseStatusOk

	return nil
}

func (s *Service) AutoCreatePayoutDocuments(
	ctx context.Context,
	_ *billingpb.EmptyRequest,
	_ *billingpb.EmptyResponse,
) error {
	zap.L().Info("start auto-creation of payout documents")

	merchants, err := s.merchantRepository.GetMerchantsWithAutoPayouts(ctx)
	if err != nil {
		zap.L().Error("GetMerchantsWithAutoPayouts failed", zap.Error(err))
		return err
	}

	req1 := &billingpb.CreatePayoutDocumentRequest{
		Ip:               "0.0.0.0",
		Initiator:        pkg.RoyaltyReportChangeSourceAuto,
		IsAutoGeneration: true,
	}
	res := &billingpb.CreatePayoutDocumentResponse{}

	wasErrors := false

	for _, m := range merchants {
		req1.MerchantId = m.Id
		err = s.createPayoutDocument(ctx, m, req1, res)

		if err != nil {
			if err == errorPayoutSourcesNotFound {
				continue
			}
			zap.L().Error(
				"auto createPayoutDocument failed with error",
				zap.Error(err),
				zap.String("merchantId", m.Id),
			)
			wasErrors = true
			continue
		}
		if res.Status != billingpb.ResponseStatusOk {
			if res.Message == errorPayoutAmountInvalid || res.Message == errorPayoutSourcesNotFound {
				continue
			}
			zap.L().Error(
				"auto createPayoutDocument failed in response",
				zap.Int32("code", res.Status),
				zap.Any("message", res.Message),
				zap.String("merchantId", m.Id),
			)
			wasErrors = true
			continue
		}
	}

	if wasErrors {
		zap.L().Warn("auto-creation of payout documents finished with errors")

		return errorPayoutAutoPayoutsWithErrors
	}

	zap.L().Info("auto-creation of payout documents finished")

	return nil
}

func (s *Service) renderPayoutDocument(
	ctx context.Context,
	pd *billingpb.PayoutDocument,
	merchant *billingpb.Merchant,
) error {
	params, err := json.Marshal(map[string]interface{}{reporterpb.ParamsFieldId: pd.Id})
	if err != nil {
		zap.L().Error(
			"Unable to marshal the params of payout for the reporting service.",
			zap.Error(err),
		)
		return err
	}

	fileReq := &reporterpb.ReportFile{
		UserId:           merchant.User.Id,
		MerchantId:       merchant.Id,
		ReportType:       reporterpb.ReportTypePayout,
		FileType:         reporterpb.OutputExtensionPdf,
		Params:           params,
		SendNotification: merchant.ManualPayoutsEnabled,
	}

	if _, err = s.reporterService.CreateFile(ctx, fileReq); err != nil {
		zap.L().Error(
			"Unable to create file in the reporting service for payout.",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *Service) UpdatePayoutDocument(
	ctx context.Context,
	req *billingpb.UpdatePayoutDocumentRequest,
	res *billingpb.PayoutDocumentResponse,
) error {
	pd, err := s.payoutRepository.GetById(ctx, req.PayoutDocumentId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPayoutNotFound
			return nil
		}
		return err
	}

	isChanged := false
	needBalanceUpdate := false

	_, isReqStatusForBecomePaid := statusForBecomePaid[req.Status]
	becomePaid := isReqStatusForBecomePaid && pd.Status != req.Status

	_, isReqStatusForBecomeFailed := statusForBecomeFailed[req.Status]
	_, isPayoutStatusForBecomeFailed := statusForBecomeFailed[pd.Status]
	becomeFailed := isReqStatusForBecomeFailed && !isPayoutStatusForBecomeFailed

	if req.Status != "" && pd.Status != req.Status {
		if pd.Status == pkg.PayoutDocumentStatusPaid || pd.Status == pkg.PayoutDocumentStatusFailed {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = errorPayoutStatusChangeIsForbidden

			return nil
		}

		if req.Status == pkg.PayoutDocumentStatusPaid {
			pd.PaidAt = ptypes.TimestampNow()
		}

		isChanged = true
		pd.Status = req.Status
		if _, ok := statusForUpdateBalance[pd.Status]; ok {
			needBalanceUpdate = true
		}
	}

	if req.Transaction != "" && pd.Transaction != req.Transaction {
		isChanged = true
		pd.Transaction = req.Transaction
	}

	if req.FailureCode != "" && pd.FailureCode != req.FailureCode {
		isChanged = true
		pd.FailureCode = req.FailureCode
	}

	if req.FailureMessage != "" && pd.FailureMessage != req.FailureMessage {
		isChanged = true
		pd.FailureMessage = req.FailureMessage
	}

	if req.FailureTransaction != "" && pd.FailureTransaction != req.FailureTransaction {
		isChanged = true
		pd.FailureTransaction = req.FailureTransaction
	}

	if isChanged {
		err = s.payoutRepository.Update(ctx, pd, req.Ip, payoutChangeSourceAdmin)
		if err != nil {
			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusSystemError
				res.Message = e
				return nil
			}
			return err
		}

		if becomePaid == true {
			err = s.royaltyReportSetPaid(ctx, pd.SourceId, pd.Id, req.Ip, pkg.RoyaltyReportChangeSourceAdmin)
			if err != nil {
				res.Status = billingpb.ResponseStatusSystemError
				res.Message = errorPayoutUpdateRoyaltyReports

				return nil
			}

		} else {
			if becomeFailed == true {
				err = s.royaltyReportUnsetPaid(ctx, pd.SourceId, req.Ip, pkg.RoyaltyReportChangeSourceAdmin)
				if err != nil {
					res.Status = billingpb.ResponseStatusSystemError
					res.Message = errorPayoutUpdateRoyaltyReports

					return nil
				}
			}
		}

		res.Status = billingpb.ResponseStatusOk
	} else {
		res.Status = billingpb.ResponseStatusNotModified
	}

	if needBalanceUpdate == true {
		_, err = s.updateMerchantBalance(ctx, pd.MerchantId)
		if err != nil {
			res.Status = billingpb.ResponseStatusSystemError
			res.Message = errorPayoutUpdateBalance

			return nil
		}
	}

	res.Item = pd
	return nil
}

func (s *Service) royaltyReportUnsetPaid(ctx context.Context, reportIds []string, ip, source string) (err error) {
	for _, id := range reportIds {
		rr, err := s.royaltyReportRepository.GetById(ctx, id)
		if err != nil {
			return err
		}

		rr.PayoutDocumentId = ""
		rr.Status = billingpb.RoyaltyReportStatusAccepted
		rr.PayoutDate = nil

		err = s.royaltyReportRepository.Update(ctx, rr, ip, source)
		if err != nil {
			return err
		}
	}
	return
}

func (s *Service) GetPayoutDocuments(
	ctx context.Context,
	req *billingpb.GetPayoutDocumentsRequest,
	res *billingpb.GetPayoutDocumentsResponse,
) error {
	res.Status = billingpb.ResponseStatusOk

	count, err := s.payoutRepository.FindCount(ctx, req.MerchantId, req.Status, req.DateFrom, req.DateTo)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if count == 0 {
		res.Status = billingpb.ResponseStatusOk
		res.Data = &billingpb.PayoutDocumentsPaginate{
			Count: 0,
			Items: nil,
		}
		return nil
	}

	pds, err := s.payoutRepository.Find(ctx, req.MerchantId, req.Status, req.DateFrom, req.DateTo, req.Offset, req.Limit)

	if err != nil {
		return err
	}

	res.Status = billingpb.ResponseStatusOk
	res.Data = &billingpb.PayoutDocumentsPaginate{
		Count: int32(count),
		Items: pds,
	}
	return nil
}

func (s *Service) PayoutDocumentPdfUploaded(
	ctx context.Context,
	req *billingpb.PayoutDocumentPdfUploadedRequest,
	res *billingpb.PayoutDocumentPdfUploadedResponse,
) error {
	pd, err := s.payoutRepository.GetById(ctx, req.PayoutId)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPayoutNotFound
			return nil
		}
		return err
	}

	merchant, err := s.merchantRepository.GetById(ctx, pd.MerchantId)

	if err != nil {
		return merchantErrorNotFound
	}

	if merchant.HasAuthorizedEmail() == false {
		zap.L().Warn("Merchant has no authorized email", zap.String("merchant_id", merchant.Id))
		res.Status = billingpb.ResponseStatusOk
		return nil
	}

	content := base64.StdEncoding.EncodeToString(req.Content)
	contentType := mime.TypeByExtension(filepath.Ext(req.Filename))

	periodFrom, err := ptypes.Timestamp(pd.PeriodFrom)
	if err != nil {
		zap.L().Error(
			pkg.ErrorTimeConversion,
			zap.Any(pkg.ErrorTimeConversionMethod, "ptypes.Timestamp"),
			zap.Any(pkg.ErrorTimeConversionValue, pd.PeriodFrom),
			zap.Error(err),
		)
		return err
	}
	periodTo, err := ptypes.Timestamp(pd.PeriodTo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorTimeConversion,
			zap.Any(pkg.ErrorTimeConversionMethod, "ptypes.Timestamp"),
			zap.Any(pkg.ErrorTimeConversionValue, pd.PeriodTo),
			zap.Error(err),
		)
		return err
	}

	operatingCompany, err := s.operatingCompanyRepository.GetById(ctx, pd.OperatingCompanyId)

	if err != nil {
		zap.L().Error("Operating company not found", zap.Error(err), zap.String("operating_company_id", pd.OperatingCompanyId))
		return err
	}

	payload := &postmarkpb.Payload{
		TemplateAlias: s.cfg.EmailTemplates.NewPayout,
		TemplateModel: map[string]string{
			"merchant_id":            merchant.Id,
			"payout_id":              pd.Id,
			"period_from":            periodFrom.Format("2006-01-02"),
			"period_to":              periodTo.Format("2006-01-02"),
			"license_agreement":      merchant.AgreementNumber,
			"status":                 pd.Status,
			"merchant_greeting":      merchant.GetOwnerName(),
			"payouts_url":            s.cfg.GetPayoutsUrl(),
			"operating_company_name": operatingCompany.Name,
		},
		To: merchant.GetOwnerEmail(),
		Attachments: []*postmarkpb.PayloadAttachment{
			{
				Name:        req.Filename,
				Content:     content,
				ContentType: contentType,
			},
		},
	}

	err = s.postmarkBroker.Publish(postmarkpb.PostmarkSenderTopicName, payload, amqp.Table{})

	if err != nil {
		zap.L().Error(
			"Publication message about merchant new payout document to queue failed",
			zap.Error(err),
			zap.Any("payout document", pd),
		)
	}

	res.Status = billingpb.ResponseStatusOk

	return nil
}

func (s *Service) getPayoutDocumentSources(
	ctx context.Context,
	merchant *billingpb.Merchant,
) ([]*billingpb.RoyaltyReport, error) {
	result, err := s.royaltyReportRepository.GetNonPayoutReports(ctx, merchant.Id, merchant.GetPayoutCurrency())

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments || result == nil || len(result) == 0 {
		return nil, errorPayoutSourcesNotFound
	}

	for _, v := range result {
		if v.Status == billingpb.RoyaltyReportStatusPending {
			return nil, errorPayoutSourcesPending
		}
		if v.Status == billingpb.RoyaltyReportStatusDispute {
			return nil, errorPayoutSourcesDispute
		}
	}

	return result, nil
}
