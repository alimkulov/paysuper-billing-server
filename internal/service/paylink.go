package service

import (
	"context"
	"fmt"
	u "github.com/PuerkitoBio/purell"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/go-querystring/query"
	"github.com/jinzhu/now"
	"github.com/paysuper/paysuper-billing-server/internal/repository"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type orderViewPaylinkStatFunc func(repository.OrderViewRepositoryInterface, context.Context, string, string, int64, int64) (*billingpb.GroupStatCommon, error)

type utmQueryParams struct {
	UtmSource   string `url:"utm_source,omitempty"`
	UtmMedium   string `url:"utm_medium,omitempty"`
	UtmCampaign string `url:"utm_campaign,omitempty"`
}

var (
	errorPaylinkExpired                      = newBillingServerErrorMsg("pl000001", "payment link expired")
	errorPaylinkNotFound                     = newBillingServerErrorMsg("pl000002", "paylink not found")
	errorPaylinkProjectMismatch              = newBillingServerErrorMsg("pl000003", "projectId mismatch for existing paylink")
	errorPaylinkExpiresInPast                = newBillingServerErrorMsg("pl000004", "paylink expiry date in past")
	errorPaylinkProductsLengthInvalid        = newBillingServerErrorMsg("pl000005", "paylink products length invalid")
	errorPaylinkProductsTypeInvalid          = newBillingServerErrorMsg("pl000006", "paylink products type invalid")
	errorPaylinkProductNotBelongToMerchant   = newBillingServerErrorMsg("pl000007", "at least one of paylink products is not belongs to merchant")
	errorPaylinkProductNotBelongToProject    = newBillingServerErrorMsg("pl000008", "at least one of paylink products is not belongs to project")
	errorPaylinkStatDataInconsistent         = newBillingServerErrorMsg("pl000009", "paylink stat data inconsistent")
	errorPaylinkProductNotFoundOrInvalidType = newBillingServerErrorMsg("pl000010", "at least one of paylink products is not found or have type differ from given products_type value")

	orderViewPaylinkStatFuncMap = map[string]orderViewPaylinkStatFunc{
		"GetPaylinkStatByCountry":  repository.OrderViewRepositoryInterface.GetPaylinkStatByCountry,
		"GetPaylinkStatByReferrer": repository.OrderViewRepositoryInterface.GetPaylinkStatByReferrer,
		"GetPaylinkStatByDate":     repository.OrderViewRepositoryInterface.GetPaylinkStatByDate,
		"GetPaylinkStatByUtm":      repository.OrderViewRepositoryInterface.GetPaylinkStatByUtm,
	}
)

// GetPaylinks returns list of all payment links
func (s *Service) GetPaylinks(
	ctx context.Context,
	req *billingpb.GetPaylinksRequest,
	res *billingpb.GetPaylinksResponse,
) error {
	count, err := s.paylinkRepository.FindCount(ctx, req.MerchantId, req.ProjectId)

	if err != nil {
		return err
	}

	res.Data = &billingpb.PaylinksPaginate{Count: int32(count)}

	if count > 0 {
		res.Data.Items, err = s.paylinkRepository.Find(ctx, req.MerchantId, req.ProjectId, req.Limit, req.Offset)
		if err != nil {
			return err
		}

		for _, pl := range res.Data.Items {
			visits, err := s.paylinkVisitsRepository.CountPaylinkVisits(ctx, pl.Id, 0, 0)
			if err == nil {
				pl.Visits = int32(visits)
			}
			pl.UpdateConversion()
			pl.IsExpired = pl.GetIsExpired()
		}
	}

	res.Status = billingpb.ResponseStatusOk

	return nil
}

// GetPaylink returns one payment link
func (s *Service) GetPaylink(
	ctx context.Context,
	req *billingpb.PaylinkRequest,
	res *billingpb.GetPaylinkResponse,
) (err error) {

	res.Item, err = s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPaylinkNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	visits, err := s.paylinkVisitsRepository.CountPaylinkVisits(ctx, res.Item.Id, 0, 0)
	if err == nil {
		res.Item.Visits = int32(visits)
	}
	res.Item.UpdateConversion()

	res.Status = billingpb.ResponseStatusOk
	return nil
}

func (s *Service) GetPaylinkTransactions(
	ctx context.Context,
	req *billingpb.GetPaylinkTransactionsRequest,
	res *billingpb.TransactionsResponse,
) error {
	pl, err := s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPaylinkNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	oid, _ := primitive.ObjectIDFromHex(pl.MerchantId)
	match := bson.M{
		"merchant_id":           oid,
		"issuer.reference_type": pkg.OrderIssuerReferenceTypePaylink,
		"issuer.reference":      pl.Id,
	}

	ts, err := s.orderViewRepository.GetTransactionsPublic(ctx, match, req.Limit, req.Offset)

	if err != nil {
		return err
	}

	res.Data = &billingpb.TransactionsPaginate{
		Count: int32(len(ts)),
		Items: ts,
	}

	return nil
}

// IncrPaylinkVisits adds a visit hit to stat
func (s *Service) IncrPaylinkVisits(
	ctx context.Context,
	req *billingpb.PaylinkRequestById,
	res *billingpb.EmptyResponse,
) error {
	err := s.paylinkVisitsRepository.IncrVisits(ctx, req.Id)
	if err != nil {
		return err
	}
	return nil
}

// GetPaylinkURL returns public url for Paylink
func (s *Service) GetPaylinkURL(
	ctx context.Context,
	req *billingpb.GetPaylinkURLRequest,
	res *billingpb.GetPaylinkUrlResponse,
) (err error) {

	res.Url, err = s.getPaylinkUrl(ctx, req.Id, req.MerchantId, req.UrlMask, req.UtmMedium, req.UtmMedium, req.UtmCampaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPaylinkNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			if err == errorPaylinkExpired {
				res.Status = billingpb.ResponseStatusGone
			} else {
				res.Status = billingpb.ResponseStatusBadData
			}

			res.Message = e
			return nil
		}
		return err
	}

	res.Status = billingpb.ResponseStatusOk
	return nil
}

// DeletePaylink deletes payment link
func (s *Service) DeletePaylink(
	ctx context.Context,
	req *billingpb.PaylinkRequest,
	res *billingpb.EmptyResponseWithStatus,
) error {
	pl, err := s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)

	if err != nil || pl.MerchantId != req.MerchantId {
		res.Status = billingpb.ResponseStatusNotFound
		res.Message = errorPaylinkNotFound
		return nil
	}

	err = s.paylinkRepository.Delete(ctx, pl)

	if err != nil {
		res.Status = billingpb.ResponseStatusNotFound
		res.Message = errorPaylinkNotFound
		return nil
	}

	res.Status = billingpb.ResponseStatusOk
	return nil
}

// CreateOrUpdatePaylink create or modify payment link
func (s *Service) CreateOrUpdatePaylink(
	ctx context.Context,
	req *billingpb.CreatePaylinkRequest,
	res *billingpb.GetPaylinkResponse,
) (err error) {
	isNew := req.GetId() == ""

	pl := &billingpb.Paylink{}

	if isNew {
		pl.Id = primitive.NewObjectID().Hex()
		pl.CreatedAt = ptypes.TimestampNow()
		pl.Object = "paylink"
		pl.MerchantId = req.MerchantId
		pl.ProjectId = req.ProjectId
		pl.ProductsType = req.ProductsType
	} else {
		pl, err = s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				res.Status = billingpb.ResponseStatusNotFound
				res.Message = errorPaylinkNotFound
				return nil
			}
			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = e
				return nil
			}
			return err
		}

		if pl.ProjectId != req.ProjectId {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = errorPaylinkProjectMismatch
			return nil
		}
	}

	pl.UpdatedAt = ptypes.TimestampNow()
	pl.Name = req.Name
	pl.NoExpiryDate = req.NoExpiryDate

	project, err := s.project.GetById(ctx, pl.ProjectId)

	if err != nil || project.MerchantId != pl.MerchantId {
		res.Status = billingpb.ResponseStatusBadData
		res.Message = projectErrorNotFound
		return nil
	}

	if pl.NoExpiryDate == false {
		expiresAt := now.New(time.Unix(req.ExpiresAt, 0)).EndOfDay()

		if time.Now().After(expiresAt) {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = errorPaylinkExpiresInPast
			return nil
		}

		pl.ExpiresAt, err = ptypes.TimestampProto(expiresAt)
		if err != nil {
			zap.L().Error(
				pkg.ErrorTimeConversion,
				zap.Any(pkg.ErrorTimeConversionMethod, "ptypes.TimestampProto"),
				zap.Any(pkg.ErrorTimeConversionValue, expiresAt),
				zap.Error(err),
			)
			return err
		}
	}

	productsLength := len(req.Products)
	if productsLength < s.cfg.PaylinkMinProducts || productsLength > s.cfg.PaylinkMaxProducts {
		res.Status = billingpb.ResponseStatusBadData
		res.Message = errorPaylinkProductsLengthInvalid
		return nil
	}

	for _, productId := range req.Products {
		switch req.ProductsType {

		case pkg.OrderType_product:
			product, err := s.productRepository.GetById(ctx, productId)
			if err != nil {
				if err.Error() == "product not found" || err == mongo.ErrNoDocuments {
					res.Status = billingpb.ResponseStatusNotFound
					res.Message = errorPaylinkProductNotFoundOrInvalidType
					return nil
				}

				if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
					res.Status = billingpb.ResponseStatusBadData
					res.Message = e
					return nil
				}
				return err
			}

			if product.MerchantId != pl.MerchantId {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = errorPaylinkProductNotBelongToMerchant
				return nil
			}

			if product.ProjectId != pl.ProjectId {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = errorPaylinkProductNotBelongToProject
				return nil
			}

			break

		case pkg.OrderType_key:
			product, err := s.keyProductRepository.GetById(ctx, productId)
			if err != nil {
				if err.Error() == "key_product not found" || err == mongo.ErrNoDocuments {
					res.Status = billingpb.ResponseStatusNotFound
					res.Message = errorPaylinkProductNotFoundOrInvalidType
					return nil
				}

				if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
					res.Status = billingpb.ResponseStatusBadData
					res.Message = e
					return nil
				}
				return err
			}

			if product.MerchantId != pl.MerchantId {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = errorPaylinkProductNotBelongToMerchant
				return nil
			}

			if product.ProjectId != pl.ProjectId {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = errorPaylinkProductNotBelongToProject
				return nil
			}
			break

		default:
			res.Status = billingpb.ResponseStatusBadData
			res.Message = errorPaylinkProductsTypeInvalid
			return nil
		}
	}

	pl.ProductsType = req.ProductsType
	pl.Products = req.Products

	if isNew {
		err = s.paylinkRepository.Insert(ctx, pl)
	} else {
		err = s.paylinkRepository.Update(ctx, pl)
	}
	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Item = pl
	res.Status = billingpb.ResponseStatusOk

	return nil
}

// GetPaylinkStatTotal returns total stat for requested paylink and period
func (s *Service) GetPaylinkStatTotal(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonResponse,
) (err error) {

	pl, err := s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPaylinkNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	visits, err := s.paylinkVisitsRepository.CountPaylinkVisits(ctx, pl.Id, req.PeriodFrom, req.PeriodTo)
	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Item, err = s.orderViewRepository.GetPaylinkStat(ctx, pl.Id, req.MerchantId, req.PeriodFrom, req.PeriodTo)
	if err != nil {
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Item.PaylinkId = pl.Id
	res.Item.Visits = int32(visits)
	res.Item.UpdateConversion()

	res.Status = billingpb.ResponseStatusOk
	return nil
}

// GetPaylinkStatByCountry returns stat groped by country for requested paylink and period
func (s *Service) GetPaylinkStatByCountry(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonGroupResponse,
) (err error) {
	err = s.getPaylinkStatGroup(ctx, req, res, "GetPaylinkStatByCountry")
	if err != nil {
		return err
	}
	return nil
}

// GetPaylinkStatByReferrer returns stat grouped by referer hosts for requested paylink and period
func (s *Service) GetPaylinkStatByReferrer(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonGroupResponse,
) (err error) {
	err = s.getPaylinkStatGroup(ctx, req, res, "GetPaylinkStatByReferrer")
	if err != nil {
		return err
	}
	return nil
}

// GetPaylinkStatByDate returns stat groped by date for requested paylink and period
func (s *Service) GetPaylinkStatByDate(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonGroupResponse,
) (err error) {
	err = s.getPaylinkStatGroup(ctx, req, res, "GetPaylinkStatByDate")
	if err != nil {
		return err
	}
	return nil
}

// GetPaylinkStatByUtm returns stat groped by utm labels for requested paylink and period
func (s *Service) GetPaylinkStatByUtm(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonGroupResponse,
) (err error) {
	err = s.getPaylinkStatGroup(ctx, req, res, "GetPaylinkStatByUtm")
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) getPaylinkStatGroup(
	ctx context.Context,
	req *billingpb.GetPaylinkStatCommonRequest,
	res *billingpb.GetPaylinkStatCommonGroupResponse,
	function string,
) (err error) {
	pl, err := s.paylinkRepository.GetByIdAndMerchant(ctx, req.Id, req.MerchantId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			res.Status = billingpb.ResponseStatusNotFound
			res.Message = errorPaylinkNotFound
			return nil
		}
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	res.Item, err = orderViewPaylinkStatFuncMap[function](s.orderViewRepository, ctx, pl.Id, pl.MerchantId, req.PeriodFrom, req.PeriodTo)
	if err != nil {
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

func (s *Service) getPaylinkUrl(ctx context.Context, id, merchantId, urlMask, utmSource, utmMedium, utmCampaign string) (string, error) {
	pl, err := s.paylinkRepository.GetByIdAndMerchant(ctx, id, merchantId)
	if err != nil {
		return "", err
	}
	if pl.GetIsExpired() {
		return "", errorPaylinkExpired
	}

	if urlMask == "" {
		urlMask = pkg.PaylinkUrlDefaultMask
	}

	urlString := fmt.Sprintf(urlMask, id)

	utmQuery := &utmQueryParams{
		UtmSource:   utmSource,
		UtmMedium:   utmMedium,
		UtmCampaign: utmCampaign,
	}

	q, err := query.Values(utmQuery)
	if err != nil {
		zap.L().Error(
			"Failed to serialize utm query params",
			zap.Error(err),
		)
		return "", err
	}
	encodedQuery := q.Encode()
	if encodedQuery != "" {
		urlString += "?" + encodedQuery
	}

	return u.NormalizeURLString(urlString, u.FlagsUsuallySafeGreedy|u.FlagRemoveDuplicateSlashes)
}
