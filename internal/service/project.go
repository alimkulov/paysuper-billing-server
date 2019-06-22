package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"go.uber.org/zap"
)

const (
	cacheProjectId = "project:id:%s"

	collectionProject = "project"

	projectErrorNotFound                  = "project with specified identifier not found"
	projectErrorNameDefaultLangRequired   = "project name in \"" + DefaultLanguage + "\" locale is required"
	projectErrorCallbackCurrencyIncorrect = "project callback currency is incorrect"
	projectErrorLimitCurrencyIncorrect    = "project limit currency is incorrect"
	projectErrorLimitCurrencyRequired     = "project limit currency can't be empty if you send min or max payment amount"
)

var (
	errProjectNotFound = errors.New(projectErrorNotFound)
)

func (s *Service) ChangeProject(
	ctx context.Context,
	req *billing.Project,
	rsp *grpc.ChangeProjectResponse,
) error {
	var project *billing.Project
	var err error

	if _, err := s.merchant.GetById(req.MerchantId); err != nil {
		rsp.Status = pkg.ResponseStatusNotFound
		rsp.Message = merchantErrorNotFound

		return nil
	}

	if req.Id != "" {
		project, err = s.getProjectBy(bson.M{"_id": bson.ObjectIdHex(req.Id), "merchant_id": bson.ObjectIdHex(req.MerchantId)})

		if err != nil {
			rsp.Status = pkg.ResponseStatusNotFound
			rsp.Message = err.Error()

			return nil
		}
	}

	if _, ok := req.Name[DefaultLanguage]; !ok {
		rsp.Status = pkg.ResponseStatusBadData
		rsp.Message = projectErrorNameDefaultLangRequired

		return nil
	}

	if req.CallbackCurrency != "" {
		if _, err := s.currency.GetByCodeA3(req.CallbackCurrency); err != nil {
			rsp.Status = pkg.ResponseStatusBadData
			rsp.Message = projectErrorCallbackCurrencyIncorrect

			return nil
		}
	}

	if req.LimitsCurrency != "" {
		if _, err := s.currency.GetByCodeA3(req.LimitsCurrency); err != nil {
			rsp.Status = pkg.ResponseStatusBadData
			rsp.Message = projectErrorLimitCurrencyIncorrect

			return nil
		}
	}

	if (req.MinPaymentAmount > 0 || req.MaxPaymentAmount > 0) && req.LimitsCurrency == "" {
		rsp.Status = pkg.ResponseStatusBadData
		rsp.Message = projectErrorLimitCurrencyRequired

		return nil
	}

	if project == nil {
		project, err = s.createProject(req)
	} else {
		err = s.updateProject(req, project)
	}

	if err != nil {
		rsp.Status = pkg.ResponseStatusSystemError
		rsp.Message = err.Error()

		return nil
	}

	rsp.Status = pkg.ResponseStatusOk
	rsp.Item = project

	return nil
}

func (s *Service) GetProject(
	ctx context.Context,
	req *grpc.GetProjectRequest,
	rsp *grpc.ChangeProjectResponse,
) error {
	query := bson.M{"_id": bson.ObjectIdHex(req.ProjectId)}

	if req.MerchantId != "" {
		query["merchant_id"] = bson.ObjectIdHex(req.MerchantId)
	}

	project, err := s.getProjectBy(query)

	if err != nil {
		rsp.Status = pkg.ResponseStatusNotFound
		rsp.Message = projectErrorNotFound

		return nil
	}

	project.ProductsCount = s.getProductsCountByProject(project.Id)

	rsp.Status = pkg.ResponseStatusOk
	rsp.Item = project

	return nil
}

func (s *Service) ListProjects(
	ctx context.Context,
	req *grpc.ListProjectsRequest,
	rsp *grpc.ListProjectsResponse,
) error {
	var projects []*billing.Project
	query := make(bson.M)

	if req.MerchantId != "" {
		query["merchant_id"] = bson.ObjectIdHex(req.MerchantId)
	}

	if req.QuickSearch != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$elemMatch": bson.M{"value": bson.RegEx{Pattern: req.QuickSearch, Options: "i"}}}},
			{"id_string": bson.RegEx{Pattern: req.QuickSearch, Options: "i"}},
		}
	}

	if len(req.Statuses) > 0 {
		query["status"] = bson.M{"$in": req.Statuses}
	}

	count, err := s.db.Collection(collectionProject).Find(query).Count()

	if err != nil {
		zap.S().Errorf("Query to count projects failed", "err", err.Error(), "query", query)
		return errors.New(orderErrorUnknown)
	}

	afQuery := []bson.M{
		{"$match": query},
		{
			"$lookup": bson.M{
				"from":         collectionProduct,
				"localField":   "_id",
				"foreignField": "project_id",
				"as":           "products",
			},
		},
		{
			"$project": bson.M{
				"_id":                         "$_id",
				"merchant_id":                 "$merchant_id",
				"name":                        "$name",
				"callback_protocol":           "$callback_protocol",
				"callback_currency":           "$callback_currency",
				"create_order_allowed_urls":   "$create_order_allowed_urls",
				"allow_dynamic_notify_urls":   "$allow_dynamic_notify_urls",
				"allow_dynamic_redirect_urls": "$allow_dynamic_redirect_urls",
				"limits_currency":             "$limits_currency",
				"min_payment_amount":          "$min_payment_amount",
				"max_payment_amount":          "$max_payment_amount",
				"notify_emails":               "$notify_emails",
				"is_products_checkout":        "$is_products_checkout",
				"secret_key":                  "$secret_key",
				"signature_required":          "$signature_required",
				"send_notify_email":           "$send_notify_email",
				"url_check_account":           "$url_check_account",
				"url_process_payment":         "$url_process_payment",
				"url_redirect_fail":           "$url_redirect_fail",
				"url_redirect_success":        "$url_redirect_success",
				"status":                      "$status",
				"created_at":                  "$created_at",
				"updated_at":                  "$updated_at",
				"products_count":              bson.M{"$size": "$products"},
			},
		},
		{"$skip": req.Offset},
		{"$limit": req.Limit},
	}

	if len(req.Sort) > 0 {
		afQuery = s.mgoPipeSort(afQuery, req.Sort)
	}

	err = s.db.Collection(collectionProject).Pipe(afQuery).All(&projects)

	if err != nil {
		zap.S().Errorf("Query to find projects failed", "err", err.Error(), "query", afQuery)
		return errors.New(orderErrorUnknown)
	}

	rsp.Count = int32(count)
	rsp.Items = []*billing.Project{}

	if count > 0 {
		rsp.Items = projects
	}

	return nil
}

func (s *Service) DeleteProject(
	ctx context.Context,
	req *grpc.GetProjectRequest,
	rsp *grpc.ChangeProjectResponse,
) error {
	query := bson.M{"_id": bson.ObjectIdHex(req.ProjectId)}

	if req.MerchantId != "" {
		query["merchant_id"] = bson.ObjectIdHex(req.MerchantId)
	}

	project, err := s.getProjectBy(query)

	if err != nil {
		rsp.Status = pkg.ResponseStatusNotFound
		rsp.Message = projectErrorNotFound

		return nil
	}

	rsp.Status = pkg.ResponseStatusOk

	if project.IsDeleted() == true {
		return nil
	}

	project.Status = pkg.ProjectStatusDeleted

	if err := s.project.Update(project); err != nil {
		zap.S().Errorf("Query to delete project failed", "err", err.Error(), "data", project)

		rsp.Status = pkg.ResponseStatusSystemError
		rsp.Message = err.Error()

		return nil
	}

	return nil
}

func (s *Service) getProjectBy(query bson.M) (project *billing.Project, err error) {
	err = s.db.Collection(collectionProject).Find(query).One(&project)

	if err != nil {
		if err != mgo.ErrNotFound {
			zap.S().Errorf("Query to find project failed", "err", err.Error(), "query", query)
		}

		return project, errProjectNotFound
	}

	return
}

func (s *Service) createProject(req *billing.Project) (*billing.Project, error) {
	project := &billing.Project{
		Id:                       bson.NewObjectId().Hex(),
		MerchantId:               req.MerchantId,
		Name:                     req.Name,
		CallbackCurrency:         req.CallbackCurrency,
		CallbackProtocol:         req.CallbackProtocol,
		CreateOrderAllowedUrls:   req.CreateOrderAllowedUrls,
		AllowDynamicNotifyUrls:   req.AllowDynamicNotifyUrls,
		AllowDynamicRedirectUrls: req.AllowDynamicRedirectUrls,
		LimitsCurrency:           req.LimitsCurrency,
		MinPaymentAmount:         req.MinPaymentAmount,
		MaxPaymentAmount:         req.MaxPaymentAmount,
		NotifyEmails:             req.NotifyEmails,
		IsProductsCheckout:       req.IsProductsCheckout,
		SecretKey:                req.SecretKey,
		SignatureRequired:        req.SignatureRequired,
		SendNotifyEmail:          req.SendNotifyEmail,
		UrlCheckAccount:          req.UrlCheckAccount,
		UrlProcessPayment:        req.UrlProcessPayment,
		UrlRedirectFail:          req.UrlRedirectFail,
		UrlRedirectSuccess:       req.UrlRedirectSuccess,
		UrlChargebackPayment:     req.UrlChargebackPayment,
		UrlCancelPayment:         req.UrlCancelPayment,
		UrlFraudPayment:          req.UrlFraudPayment,
		UrlRefundPayment:         req.UrlRefundPayment,
		Status:                   pkg.ProjectStatusDraft,
		CreatedAt:                ptypes.TimestampNow(),
		UpdatedAt:                ptypes.TimestampNow(),
	}

	if err := s.project.Insert(project); err != nil {
		zap.S().Errorf("Query to create project failed", "err", err.Error(), "data", project)
		return nil, errors.New(orderErrorUnknown)
	}

	return project, nil
}

func (s *Service) updateProject(req *billing.Project, project *billing.Project) error {
	project.Name = req.Name
	project.CallbackCurrency = req.CallbackCurrency
	project.CreateOrderAllowedUrls = req.CreateOrderAllowedUrls
	project.AllowDynamicNotifyUrls = req.AllowDynamicNotifyUrls
	project.AllowDynamicRedirectUrls = req.AllowDynamicRedirectUrls
	project.LimitsCurrency = req.LimitsCurrency
	project.MinPaymentAmount = req.MinPaymentAmount
	project.MaxPaymentAmount = req.MaxPaymentAmount
	project.NotifyEmails = req.NotifyEmails
	project.IsProductsCheckout = req.IsProductsCheckout
	project.SecretKey = req.SecretKey
	project.SignatureRequired = req.SignatureRequired
	project.SendNotifyEmail = req.SendNotifyEmail
	project.UrlRedirectFail = req.UrlRedirectFail
	project.UrlRedirectSuccess = req.UrlRedirectSuccess
	project.Status = req.Status
	project.UpdatedAt = ptypes.TimestampNow()

	if project.NeedChangeStatusToDraft(req) == true {
		project.Status = pkg.ProjectStatusDraft
	}

	project.CallbackProtocol = req.CallbackProtocol
	project.UrlCheckAccount = req.UrlCheckAccount
	project.UrlProcessPayment = req.UrlProcessPayment
	project.UrlChargebackPayment = req.UrlChargebackPayment
	project.UrlCancelPayment = req.UrlCancelPayment
	project.UrlFraudPayment = req.UrlFraudPayment
	project.UrlRefundPayment = req.UrlRefundPayment

	if err := s.project.Update(project); err != nil {
		zap.S().Errorf("Query to update project failed", "err", err.Error(), "data", project)
		return errors.New(orderErrorUnknown)
	}

	project.ProductsCount = s.getProductsCountByProject(project.Id)

	return nil
}

func newProjectService(svc *Service) *Project {
	s := &Project{svc: svc}
	return s
}

func (h *Project) Insert(project *billing.Project) error {
	if err := h.svc.db.Collection(collectionProject).Insert(project); err != nil {
		return err
	}

	if err := h.svc.cacher.Set(fmt.Sprintf(cacheProjectId, project.Id), project, 0); err != nil {
		return err
	}

	return nil
}

func (h *Project) MultipleInsert(projects []*billing.Project) error {
	p := make([]interface{}, len(projects))
	for i, v := range projects {
		p[i] = v
	}

	if err := h.svc.db.Collection(collectionProject).Insert(p...); err != nil {
		return err
	}

	return nil
}

func (h *Project) Update(project *billing.Project) error {
	if err := h.svc.db.Collection(collectionProject).UpdateId(bson.ObjectIdHex(project.Id), project); err != nil {
		return err
	}

	if err := h.svc.cacher.Set(fmt.Sprintf(cacheProjectId, project.Id), project, 0); err != nil {
		return err
	}

	return nil
}

func (h Project) GetById(id string) (*billing.Project, error) {
	var c billing.Project
	key := fmt.Sprintf(cacheProjectId, id)

	if err := h.svc.cacher.Get(key, c); err == nil {
		return &c, nil
	}

	if err := h.svc.db.Collection(collectionProject).
		Find(bson.M{"_id": bson.ObjectIdHex(id)}).
		One(&c); err != nil {
		return nil, fmt.Errorf(errorNotFound, collectionProject)
	}

	_ = h.svc.cacher.Set(key, c, 0)
	return &c, nil
}
