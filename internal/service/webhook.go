package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	constant "github.com/paysuper/paysuper-proto/go/recurringpb"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	webhookTypeIncorrect = newBillingServerErrorMsg("wh000001", "type for webhook is invalid")
)

const (
	orderRequestType = "type"
)

func (s *Service) SendWebhookToMerchant(
	ctx context.Context,
	req *billingpb.OrderCreateRequest,
	res *billingpb.SendWebhookToMerchantResponse,
) error {
	res.Status = billingpb.ResponseStatusOk

	processor := &OrderCreateRequestProcessor{
		Service: s,
		request: req,
		checked: &orderCreateRequestProcessorChecked{},
	}

	if err := processor.processProject(); err != nil {
		zap.S().Errorw(pkg.MethodFinishedWithError, "err", err.Error())
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	switch req.Type {
	case pkg.OrderType_key:
		if err := processor.processPaylinkKeyProducts(); err != nil {
			if pid := req.PrivateMetadata["PaylinkId"]; pid != "" {
				s.notifyPaylinkError(ctx, pid, err, req, nil)
			}
			zap.S().Errorw(pkg.MethodFinishedWithError, "err", err.Error())
			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = e
				return nil
			}
			return err
		}
		break
	case pkg.OrderType_product:
		if err := processor.processPaylinkProducts(ctx); err != nil {
			if pid := req.PrivateMetadata["PaylinkId"]; pid != "" {
				s.notifyPaylinkError(ctx, pid, err, req, nil)
			}

			if err == billingpb.ProductNoPriceInCurrencyError {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = productNoPriceInCurrencyError
				return nil
			}

			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = e
				return nil
			}
			return err
		}
		break
	case pkg.OrderTypeVirtualCurrency:
		err := processor.processVirtualCurrency(ctx)
		if err != nil {
			zap.L().Error(
				pkg.MethodFinishedWithError,
				zap.Error(err),
			)

			res.Status = billingpb.ResponseStatusBadData
			res.Message = err.(*billingpb.ResponseErrorMessage)
			return nil
		}
		break
	case pkg.OrderType_simple:
		processor.processAmount()
	default:
		zap.L().Error(
			webhookTypeIncorrect.Message,
			zap.String(orderRequestType, req.Type),
		)
		res.Status = billingpb.ResponseStatusBadData
		res.Message = webhookTypeIncorrect
		res.Message.Details = req.Type
		return nil
	}

	if req.User != nil {
		err := processor.processUserData()

		if err != nil {
			zap.S().Errorw(pkg.MethodFinishedWithError, "err", err.Error())
			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = e
				return nil
			}
			return err
		}
	}

	err := processor.processCurrency(req.Type)
	if err != nil {
		zap.L().Error("process currency failed", zap.Error(err))
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	if processor.checked.currency == "" {
		res.Status = billingpb.ResponseStatusBadData
		res.Message = orderErrorCurrencyIsRequired
		return nil
	}

	if req.OrderId != "" {
		if err := processor.processProjectOrderId(); err != nil {
			zap.S().Errorw(pkg.MethodFinishedWithError, "err", err.Error())
			if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
				res.Status = billingpb.ResponseStatusBadData
				res.Message = e
				return nil
			}
			return err
		}
	}

	processor.processMetadata()
	processor.processPrivateMetadata()

	order, err := processor.prepareOrder()

	if err != nil {
		zap.S().Errorw(pkg.MethodFinishedWithError, "err", err.Error())
		if e, ok := err.(*billingpb.ResponseErrorMessage); ok {
			res.Status = billingpb.ResponseStatusBadData
			res.Message = e
			return nil
		}
		return err
	}

	order.PrivateStatus = constant.OrderStatusPaymentSystemComplete

	if req.TestingCase == billingpb.TestCaseNonExistingUser {
		order.User.ExternalId = "paysuper_test_" + uuid.New().String()
	}

	err = s.broker.Publish(constant.PayOneTopicNotifyPaymentName, order, amqp.Table{"x-retry-count": int32(0)})

	if err != nil {
		zap.L().Error(
			orderErrorPublishNotificationFailed,
			zap.Error(err),
			zap.String("topic", constant.PayOneTopicNotifyPaymentName),
			zap.Any("order", order),
		)
	}

	res.OrderId = order.GetUuid()

	return nil
}

func (s *Service) NotifyWebhookTestResults(
	ctx context.Context,
	req *billingpb.NotifyWebhookTestResultsRequest,
	res *billingpb.EmptyResponseWithStatus,
) error {
	res.Status = billingpb.ResponseStatusOk

	project, err := s.project.GetById(ctx, req.ProjectId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorGrpcServiceCallFailed,
			zap.Error(err),
			zap.String(errorFieldService, "Project"),
			zap.String(errorFieldMethod, "GetById"),
		)
		res.Status = billingpb.ResponseStatusSystemError
		res.Message = projectErrorUnknown
		return nil
	}

	if project.WebhookTesting == nil {
		project.WebhookTesting = &billingpb.WebHookTesting{}
	}

	switch req.Type {
	case pkg.OrderType_product:
		s.processTestingProducts(project, req)
		break
	case pkg.OrderType_key:
		s.processTestingKeys(project, req)
		break
	case pkg.OrderTypeVirtualCurrency:
		s.processTestingVirtualCurrency(project, req)
		break
	default:
		zap.L().Error(
			pkg.UnknownTypeError,
			zap.Error(err),
			zap.String(errorFieldService, "Project"),
			zap.String(errorFieldMethod, "GetById"),
		)
		res.Status = billingpb.ResponseStatusBadData
		res.Message = webhookTypeIncorrect
		res.Message.Details = req.Type
		return nil
	}

	err = s.project.Update(ctx, project)

	if err != nil {
		res.Status = billingpb.ResponseStatusSystemError
		res.Message = projectErrorUnknown
		return nil
	}

	return nil
}

func (s *Service) processTestingVirtualCurrency(project *billingpb.Project, req *billingpb.NotifyWebhookTestResultsRequest) {
	if project.WebhookTesting.VirtualCurrency == nil {
		project.WebhookTesting.VirtualCurrency = &billingpb.VirtualCurrencyTesting{}
	}
	switch req.TestCase {
	case pkg.TestCaseNonExistingUser:
		project.WebhookTesting.VirtualCurrency.NonExistingUser = req.IsPassed
		break
	case pkg.TestCaseExistingUser:
		project.WebhookTesting.VirtualCurrency.ExistingUser = req.IsPassed
		break
	case pkg.TestCaseCorrectPayment:
		project.WebhookTesting.VirtualCurrency.CorrectPayment = req.IsPassed
		break
	case pkg.TestCaseIncorrectPayment:
		project.WebhookTesting.VirtualCurrency.IncorrectPayment = req.IsPassed
		break
	}
}

func (s *Service) processTestingKeys(project *billingpb.Project, req *billingpb.NotifyWebhookTestResultsRequest) {
	if project.WebhookTesting.Keys == nil {
		project.WebhookTesting.Keys = &billingpb.KeysTesting{}
	}
	project.WebhookTesting.Keys.IsPassed = req.IsPassed
}

func (s *Service) processTestingProducts(project *billingpb.Project, req *billingpb.NotifyWebhookTestResultsRequest) {
	if project.WebhookTesting.Products == nil {
		project.WebhookTesting.Products = &billingpb.ProductsTesting{}
	}
	switch req.TestCase {
	case pkg.TestCaseNonExistingUser:
		project.WebhookTesting.Products.NonExistingUser = req.IsPassed
		break
	case pkg.TestCaseExistingUser:
		project.WebhookTesting.Products.ExistingUser = req.IsPassed
		break
	case pkg.TestCaseCorrectPayment:
		project.WebhookTesting.Products.CorrectPayment = req.IsPassed
		break
	case pkg.TestCaseIncorrectPayment:
		project.WebhookTesting.Products.IncorrectPayment = req.IsPassed
		break
	}
}
