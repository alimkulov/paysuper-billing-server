package service

import (
	"context"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/paysuper/paysuper-billing-server/internal/config"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	"github.com/paysuper/paysuper-billing-server/internal/mock"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type ProjectCRUDTestSuite struct {
	suite.Suite
	service *Service
	cache   CacheInterface

	merchant *billing.Merchant
	project  *billing.Project
}

func Test_ProjectCRUD(t *testing.T) {
	suite.Run(t, new(ProjectCRUDTestSuite))
}

func (suite *ProjectCRUDTestSuite) SetupTest() {
	cfg, err := config.NewConfig()
	assert.NoError(suite.T(), err, "Config load failed")

	cfg.AccountingCurrency = "RUB"

	settings := database.Connection{
		Host:     cfg.MongoHost,
		Database: cfg.MongoDatabase,
		User:     cfg.MongoUser,
		Password: cfg.MongoPassword,
	}

	db, err := database.NewDatabase(settings)
	assert.NoError(suite.T(), err, "Database connection failed")
	projectId := bson.NewObjectId().Hex()

	rub := &billing.Currency{
		CodeInt:  643,
		CodeA3:   "RUB",
		Name:     &billing.Name{Ru: "Российский рубль", En: "Russian ruble"},
		IsActive: true,
	}

	pm1 := &billing.PaymentMethod{
		Id:               bson.NewObjectId().Hex(),
		Name:             "Bank card",
		Group:            "BANKCARD",
		MinPaymentAmount: 100,
		MaxPaymentAmount: 15000,
		Currency:         rub,
		Currencies:       []int32{643, 840, 980},
		Params: &billing.PaymentMethodParams{
			Handler:          "cardpay",
			Terminal:         "15985",
			Password:         "A1tph4I6BD0f",
			CallbackPassword: "0V1rJ7t4jCRv",
			ExternalId:       "BANKCARD",
		},
		Type:          "bank_card",
		IsActive:      true,
		AccountRegexp: "^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$",
		PaymentSystem: &billing.PaymentSystem{
			Id:                 bson.NewObjectId().Hex(),
			Name:               "CardPay",
			AccountingCurrency: rub,
			AccountingPeriod:   "every-day",
			Country:            &billing.Country{},
			IsActive:           true,
		},
	}

	pm2 := &billing.PaymentMethod{
		Id:               bson.NewObjectId().Hex(),
		Name:             "Bitcoin",
		Group:            "BITCOIN_1",
		MinPaymentAmount: 0,
		MaxPaymentAmount: 0,
		Currency:         rub,
		Currencies:       []int32{643, 840, 980},
		Params: &billing.PaymentMethodParams{
			Handler:    "unit_test",
			Terminal:   "16007",
			ExternalId: "BITCOIN",
		},
		Type:     "crypto",
		IsActive: true,
		PaymentSystem: &billing.PaymentSystem{
			Id:                 bson.NewObjectId().Hex(),
			Name:               "CardPay",
			AccountingCurrency: rub,
			AccountingPeriod:   "every-day",
			Country:            &billing.Country{},
			IsActive:           true,
		},
	}

	err = db.Collection(pkg.CollectionCurrency).Insert(rub)
	assert.NoError(suite.T(), err, "Insert currency test data failed")

	paymentMethods := map[string]*billing.MerchantPaymentMethod{
		pm1.Id: {
			Commission: &billing.MerchantPaymentMethodCommissions{
				Fee:            1,
				PerTransaction: &billing.MerchantPaymentMethodPerTransactionCommission{},
			},
			PaymentMethod: &billing.MerchantPaymentMethodIdentification{
				Id: pm1.Id,
			},
		},
		pm2.Id: {
			Commission: &billing.MerchantPaymentMethodCommissions{
				Fee:            1,
				PerTransaction: &billing.MerchantPaymentMethodPerTransactionCommission{},
			},
			PaymentMethod: &billing.MerchantPaymentMethodIdentification{
				Id: pm2.Id,
			},
		},
	}
	merchant := &billing.Merchant{
		Id: bson.NewObjectId().Hex(),
		User: &billing.MerchantUser{
			Id:    bson.NewObjectId().Hex(),
			Email: "test@unit.test",
		},
		Name: "Unit test",
		Zip:  "190000",
		City: "St.Petersburg",
		Contacts: &billing.MerchantContact{
			Authorized: &billing.MerchantContactAuthorized{},
			Technical:  &billing.MerchantContactTechnical{},
		},
		Banking: &billing.MerchantBanking{
			Currency: rub,
			Name:     "Bank name",
		},
		IsVatEnabled:              true,
		IsCommissionToUserEnabled: true,
		Status:                    pkg.MerchantStatusDraft,
		IsSigned:                  true,
		PaymentMethods:            paymentMethods,
	}

	project := &billing.Project{
		Id:                       projectId,
		CallbackCurrency:         "RUB",
		CallbackProtocol:         pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:           "RUB",
		MaxPaymentAmount:         15000,
		MinPaymentAmount:         1,
		Name:                     map[string]string{"en": "test project 1"},
		IsProductsCheckout:       false,
		AllowDynamicRedirectUrls: true,
		SecretKey:                "test project 1 secret key",
		Status:                   pkg.ProjectStatusInProduction,
		MerchantId:               merchant.Id,
	}

	err = db.Collection(pkg.CollectionProject).Insert(project)
	assert.NoError(suite.T(), err, "Insert project test data failed")

	products := []interface{}{
		&grpc.Product{
			Object:          "product",
			Type:            "simple_product",
			Sku:             "ru_double_yeti",
			Name:            map[string]string{"en": initialName},
			DefaultCurrency: "USD",
			Enabled:         true,
			Description:     map[string]string{"en": "blah-blah-blah"},
			LongDescription: map[string]string{"en": "Super game steam keys"},
			Url:             "http://test.ru/dffdsfsfs",
			Images:          []string{"/home/image.jpg"},
			MerchantId:      merchant.Id,
			ProjectId:       project.Id,
			Metadata: map[string]string{
				"SomeKey": "SomeValue",
			},
			Prices: []*grpc.ProductPrice{{Currency: "USD", Amount: 1005.00}},
		},
		&grpc.Product{
			Object:          "product1",
			Type:            "simple_product",
			Sku:             "ru_double_yeti1",
			Name:            map[string]string{"en": initialName},
			DefaultCurrency: "USD",
			Enabled:         true,
			Description:     map[string]string{"en": "blah-blah-blah"},
			LongDescription: map[string]string{"en": "Super game steam keys"},
			Url:             "http://test.ru/dffdsfsfs",
			Images:          []string{"/home/image.jpg"},
			MerchantId:      merchant.Id,
			ProjectId:       project.Id,
			Metadata: map[string]string{
				"SomeKey": "SomeValue",
			},
			Prices: []*grpc.ProductPrice{{Currency: "USD", Amount: 1005.00}},
		},
		&grpc.Product{
			Object:          "product2",
			Type:            "simple_product",
			Sku:             "ru_double_yeti2",
			Name:            map[string]string{"en": initialName},
			DefaultCurrency: "USD",
			Enabled:         true,
			Description:     map[string]string{"en": "blah-blah-blah"},
			LongDescription: map[string]string{"en": "Super game steam keys"},
			Url:             "http://test.ru/dffdsfsfs",
			Images:          []string{"/home/image.jpg"},
			MerchantId:      merchant.Id,
			ProjectId:       project.Id,
			Metadata: map[string]string{
				"SomeKey": "SomeValue",
			},
			Prices: []*grpc.ProductPrice{{Currency: "USD", Amount: 1005.00}},
		},
	}

	err = db.Collection(pkg.CollectionProduct).Insert(products...)
	assert.NoError(suite.T(), err, "Insert product test data failed")

	redisdb := mock.NewTestRedis()
	suite.cache = NewCacheRedis(redisdb)
	suite.service = NewBillingService(db, cfg, make(chan bool, 1), nil, nil, nil, nil, nil, suite.cache)

	if err := suite.service.Init(); err != nil {
		suite.FailNow("Billing service initialization failed", "%v", err)
	}

	pms := []*billing.PaymentMethod{pm1, pm2}
	if err := suite.service.paymentMethod.MultipleInsert(pms); err != nil {
		suite.FailNow("Insert payment methods test data failed", "%v", err)
	}

	if err := suite.service.merchant.Insert(merchant); err != nil {
		suite.FailNow("Insert merchant test data failed", "%v", err)
	}

	suite.merchant = merchant
	suite.project = project
}

func (suite *ProjectCRUDTestSuite) TearDownTest() {
	if err := suite.service.db.Drop(); err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	suite.service.db.Close()
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_NewProject_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Empty(suite.T(), rsp.Message)
	assert.NotNil(suite.T(), rsp.Item)

	assert.NotEmpty(suite.T(), rsp.Item.Id)
	assert.Equal(suite.T(), req.MerchantId, rsp.Item.MerchantId)
	assert.Equal(suite.T(), req.Name, rsp.Item.Name)
	assert.Equal(suite.T(), req.CallbackCurrency, rsp.Item.CallbackCurrency)
	assert.Equal(suite.T(), req.CallbackProtocol, rsp.Item.CallbackProtocol)
	assert.Equal(suite.T(), req.LimitsCurrency, rsp.Item.LimitsCurrency)
	assert.Equal(suite.T(), req.MinPaymentAmount, rsp.Item.MinPaymentAmount)
	assert.Equal(suite.T(), req.MaxPaymentAmount, rsp.Item.MaxPaymentAmount)
	assert.Equal(suite.T(), req.IsProductsCheckout, rsp.Item.IsProductsCheckout)
	assert.Equal(suite.T(), pkg.ProjectStatusDraft, rsp.Item.Status)
	assert.Equal(suite.T(), int32(0), rsp.Item.ProductsCount)

	project, err := suite.service.getProjectBy(bson.M{"_id": bson.ObjectIdHex(rsp.Item.Id)})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), project)

	assert.Equal(suite.T(), project.Id, rsp.Item.Id)
	assert.Equal(suite.T(), project.MerchantId, rsp.Item.MerchantId)
	assert.Equal(suite.T(), project.Name, rsp.Item.Name)
	assert.Equal(suite.T(), project.CallbackCurrency, rsp.Item.CallbackCurrency)
	assert.Equal(suite.T(), project.CallbackProtocol, rsp.Item.CallbackProtocol)
	assert.Equal(suite.T(), project.LimitsCurrency, rsp.Item.LimitsCurrency)
	assert.Equal(suite.T(), project.MinPaymentAmount, rsp.Item.MinPaymentAmount)
	assert.Equal(suite.T(), project.MaxPaymentAmount, rsp.Item.MaxPaymentAmount)
	assert.Equal(suite.T(), project.IsProductsCheckout, rsp.Item.IsProductsCheckout)
	assert.Equal(suite.T(), project.Status, rsp.Item.Status)

	cProject, err := suite.service.project.GetById(project.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), project.Id, cProject.Id)
	assert.Equal(suite.T(), project.MerchantId, cProject.MerchantId)
	assert.Equal(suite.T(), project.Name, cProject.Name)
	assert.Equal(suite.T(), project.CallbackCurrency, cProject.CallbackCurrency)
	assert.Equal(suite.T(), project.CallbackProtocol, cProject.CallbackProtocol)
	assert.Equal(suite.T(), project.LimitsCurrency, cProject.LimitsCurrency)
	assert.Equal(suite.T(), project.MinPaymentAmount, cProject.MinPaymentAmount)
	assert.Equal(suite.T(), project.MaxPaymentAmount, cProject.MaxPaymentAmount)
	assert.Equal(suite.T(), project.IsProductsCheckout, cProject.IsProductsCheckout)
	assert.Equal(suite.T(), project.Status, cProject.Status)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_ExistProject_Ok() {
	req := suite.project
	req.Name["ua"] = "модульний тест"
	req.CallbackProtocol = pkg.ProjectCallbackProtocolDefault
	req.SecretKey = "qwerty"

	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Empty(suite.T(), rsp.Message)

	assert.Equal(suite.T(), req.Id, rsp.Item.Id)
	assert.Equal(suite.T(), req.MerchantId, rsp.Item.MerchantId)
	assert.Equal(suite.T(), req.Name, rsp.Item.Name)
	assert.Equal(suite.T(), req.CallbackProtocol, rsp.Item.CallbackProtocol)
	assert.NotEqual(suite.T(), req.Status, rsp.Item.Status)
	assert.Equal(suite.T(), pkg.ProjectStatusDraft, rsp.Item.Status)

	project, err := suite.service.getProjectBy(bson.M{"_id": bson.ObjectIdHex(rsp.Item.Id)})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), project)

	assert.Equal(suite.T(), project.Id, rsp.Item.Id)
	assert.Equal(suite.T(), project.MerchantId, rsp.Item.MerchantId)
	assert.Equal(suite.T(), project.Name, rsp.Item.Name)
	assert.Equal(suite.T(), project.CallbackCurrency, rsp.Item.CallbackCurrency)
	assert.Equal(suite.T(), project.CallbackProtocol, rsp.Item.CallbackProtocol)
	assert.Equal(suite.T(), project.LimitsCurrency, rsp.Item.LimitsCurrency)
	assert.Equal(suite.T(), project.MinPaymentAmount, rsp.Item.MinPaymentAmount)
	assert.Equal(suite.T(), project.MaxPaymentAmount, rsp.Item.MaxPaymentAmount)
	assert.Equal(suite.T(), project.IsProductsCheckout, rsp.Item.IsProductsCheckout)
	assert.Equal(suite.T(), project.Status, rsp.Item.Status)

	cProject, err := suite.service.project.GetById(project.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), project.Id, cProject.Id)
	assert.Equal(suite.T(), project.MerchantId, cProject.MerchantId)
	assert.Equal(suite.T(), project.Name, cProject.Name)
	assert.Equal(suite.T(), project.CallbackCurrency, cProject.CallbackCurrency)
	assert.Equal(suite.T(), project.CallbackProtocol, cProject.CallbackProtocol)
	assert.Equal(suite.T(), project.LimitsCurrency, cProject.LimitsCurrency)
	assert.Equal(suite.T(), project.MinPaymentAmount, cProject.MinPaymentAmount)
	assert.Equal(suite.T(), project.MaxPaymentAmount, cProject.MaxPaymentAmount)
	assert.Equal(suite.T(), project.IsProductsCheckout, cProject.IsProductsCheckout)
	assert.Equal(suite.T(), project.Status, cProject.Status)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_MerchantNotFound_Error() {
	req := &billing.Project{
		MerchantId:         bson.NewObjectId().Hex(),
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusNotFound, rsp.Status)
	assert.Equal(suite.T(), merchantErrorNotFound, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_ExistProjectIdNotFound_Error() {
	req := &billing.Project{
		Id:                 bson.NewObjectId().Hex(),
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusNotFound, rsp.Status)
	assert.Equal(suite.T(), projectErrorNotFound, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_NameInDefaultLanguageNotSet_Error() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusBadData, rsp.Status)
	assert.Equal(suite.T(), projectErrorNameDefaultLangRequired, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_CallbackCurrencyNotFound_Error() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "USD",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusBadData, rsp.Status)
	assert.Equal(suite.T(), projectErrorCallbackCurrencyIncorrect, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ChangeProject_LimitCurrencyNotFound_Error() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "USD",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusBadData, rsp.Status)
	assert.Equal(suite.T(), projectErrorLimitCurrencyIncorrect, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_GetProject_Ok() {
	req := &grpc.GetProjectRequest{
		ProjectId:  suite.project.Id,
		MerchantId: suite.merchant.Id,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.GetProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Empty(suite.T(), rsp.Message)
	assert.NotNil(suite.T(), rsp.Item)

	assert.Equal(suite.T(), suite.project.Id, rsp.Item.Id)
	assert.Equal(suite.T(), suite.project.MerchantId, rsp.Item.MerchantId)
	assert.Equal(suite.T(), suite.project.Name, rsp.Item.Name)
	assert.Equal(suite.T(), suite.project.CallbackCurrency, rsp.Item.CallbackCurrency)
	assert.Equal(suite.T(), suite.project.CallbackProtocol, rsp.Item.CallbackProtocol)
	assert.Equal(suite.T(), suite.project.LimitsCurrency, rsp.Item.LimitsCurrency)
	assert.Equal(suite.T(), suite.project.MinPaymentAmount, rsp.Item.MinPaymentAmount)
	assert.Equal(suite.T(), suite.project.MaxPaymentAmount, rsp.Item.MaxPaymentAmount)
	assert.Equal(suite.T(), suite.project.IsProductsCheckout, rsp.Item.IsProductsCheckout)
	assert.Equal(suite.T(), suite.project.Status, rsp.Item.Status)
	assert.Equal(suite.T(), int32(3), rsp.Item.ProductsCount)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_GetProject_NotFound_Error() {
	req := &grpc.GetProjectRequest{
		ProjectId:  suite.project.Id,
		MerchantId: bson.NewObjectId().Hex(),
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.GetProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusNotFound, rsp.Status)
	assert.Equal(suite.T(), projectErrorNotFound, rsp.Message)
	assert.Nil(suite.T(), rsp.Item)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ListProjects_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit1 test", "ru": "Юнит1 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit11 test", "ru": "Юнит11 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit2 test", "ru": "Юнит2 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req1 := &grpc.ListProjectsRequest{
		MerchantId: suite.merchant.Id,
		Limit:      100,
	}
	rsp1 := &grpc.ListProjectsResponse{}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(5), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 5)
	assert.Equal(suite.T(), int32(3), rsp1.Items[0].ProductsCount)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ListProjects_NameQuery_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit1 test", "ru": "Юнит1 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit11 test", "ru": "Юнит11 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit2 test", "ru": "Юнит2 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req1 := &grpc.ListProjectsRequest{
		MerchantId:  suite.merchant.Id,
		QuickSearch: "nit1",
		Limit:       100,
	}
	rsp1 := &grpc.ListProjectsResponse{}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(2), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 2)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ListProjects_StatusQuery_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "Unit test", "ru": "Юнит тест"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	rsp.Item.Status = pkg.ProjectStatusTestCompleted
	err = suite.service.ChangeProject(context.TODO(), rsp.Item, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit1 test", "ru": "Юнит1 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	rsp.Item.Status = pkg.ProjectStatusTestCompleted
	err = suite.service.ChangeProject(context.TODO(), rsp.Item, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit11 test", "ru": "Юнит11 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	rsp.Item.Status = pkg.ProjectStatusInProduction
	err = suite.service.ChangeProject(context.TODO(), rsp.Item, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "Unit2 test", "ru": "Юнит2 тест"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req1 := &grpc.ListProjectsRequest{
		MerchantId: suite.merchant.Id,
		Statuses:   []int32{pkg.ProjectStatusInProduction},
		Limit:      100,
	}
	rsp1 := &grpc.ListProjectsResponse{}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(2), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 2)

	req1.Statuses = []int32{pkg.ProjectStatusTestCompleted}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(2), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 2)

	req1.Statuses = []int32{pkg.ProjectStatusDraft, pkg.ProjectStatusTestCompleted}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(3), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 3)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_ListProjects_SortQuery_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "B", "ru": "Б"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "C", "ru": "В"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req.Name = map[string]string{"en": "D", "ru": "Г"}
	err = suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)

	req1 := &grpc.ListProjectsRequest{
		MerchantId: suite.merchant.Id,
		Sort:       []string{"name"},
		Limit:      100,
	}
	rsp1 := &grpc.ListProjectsResponse{}
	err = suite.service.ListProjects(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int32(5), rsp1.Count)
	assert.Len(suite.T(), rsp1.Items, 5)
	assert.Equal(suite.T(), "A", rsp1.Items[0].Name["en"])
	assert.Equal(suite.T(), "А", rsp1.Items[0].Name["ru"])
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_DeleteProject_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Equal(suite.T(), pkg.ProjectStatusDraft, rsp.Item.Status)

	req1 := &grpc.GetProjectRequest{
		MerchantId: req.MerchantId,
		ProjectId:  rsp.Item.Id,
	}
	rsp1 := &grpc.ChangeProjectResponse{}
	err = suite.service.DeleteProject(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp1.Status)

	project, err := suite.service.getProjectBy(bson.M{"_id": bson.ObjectIdHex(rsp.Item.Id)})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ProjectStatusDeleted, project.Status)

	project1, err := suite.service.project.GetById(rsp.Item.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), project.Status, project1.Status)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_DeleteProject_NotFound_Error() {
	req := &grpc.GetProjectRequest{
		MerchantId: suite.merchant.Id,
		ProjectId:  bson.NewObjectId().Hex(),
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.DeleteProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusNotFound, rsp.Status)
	assert.Equal(suite.T(), projectErrorNotFound, rsp.Message)
}

func (suite *ProjectCRUDTestSuite) TestProjectCRUD_DeleteDeletedProject_Ok() {
	req := &billing.Project{
		MerchantId:         suite.merchant.Id,
		Name:               map[string]string{"en": "A", "ru": "А"},
		CallbackCurrency:   "RUB",
		CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
		LimitsCurrency:     "RUB",
		MinPaymentAmount:   0,
		MaxPaymentAmount:   15000,
		IsProductsCheckout: false,
	}
	rsp := &grpc.ChangeProjectResponse{}
	err := suite.service.ChangeProject(context.TODO(), req, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Equal(suite.T(), pkg.ProjectStatusDraft, rsp.Item.Status)

	rsp.Item.Status = pkg.ProjectStatusDeleted
	err = suite.service.ChangeProject(context.TODO(), rsp.Item, rsp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp.Status)
	assert.Equal(suite.T(), pkg.ProjectStatusDeleted, rsp.Item.Status)

	req1 := &grpc.GetProjectRequest{
		MerchantId: req.MerchantId,
		ProjectId:  rsp.Item.Id,
	}
	rsp1 := &grpc.ChangeProjectResponse{}
	err = suite.service.DeleteProject(context.TODO(), req1, rsp1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), pkg.ResponseStatusOk, rsp1.Status)
}

type ProjectTestSuite struct {
	suite.Suite
	service *Service
	cache   CacheInterface
	log     *zap.Logger
	project *billing.Project
}

func Test_Project(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}

func (suite *ProjectTestSuite) SetupTest() {
	cfg, err := config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", "%v", err)
	}
	cfg.AccountingCurrency = "RUB"

	settings := database.Connection{
		Host:     cfg.MongoHost,
		Database: cfg.MongoDatabase,
		User:     cfg.MongoUser,
		Password: cfg.MongoPassword,
	}

	db, err := database.NewDatabase(settings)

	if err != nil {
		suite.FailNow("Database connection failed", "%v", err)
	}

	rub := &billing.Currency{
		CodeInt:  643,
		CodeA3:   "RUB",
		Name:     &billing.Name{Ru: "Российский рубль", En: "Russian ruble"},
		IsActive: true,
	}
	currency := []interface{}{rub}
	err = db.Collection(pkg.CollectionCurrency).Insert(currency...)
	if err != nil {
		suite.FailNow("Insert currency test data failed", "%v", err)
	}

	suite.project = &billing.Project{
		Id:                 bson.NewObjectId().Hex(),
		MerchantId:         bson.NewObjectId().Hex(),
		CallbackCurrency:   rub.CodeA3,
		CallbackProtocol:   "default",
		LimitsCurrency:     rub.CodeA3,
		MaxPaymentAmount:   15000,
		MinPaymentAmount:   0,
		Name:               map[string]string{"en": "test project 1"},
		IsProductsCheckout: true,
		SecretKey:          "test project 1 secret key",
		Status:             pkg.ProjectStatusInProduction,
	}
	err = db.Collection(pkg.CollectionProject).Insert(suite.project)
	if err != nil {
		suite.FailNow("Insert project test data failed", "%v", err)
	}

	suite.log, err = zap.NewProduction()

	if err != nil {
		suite.FailNow("Logger initialization failed", "%v", err)
	}

	redisdb := mock.NewTestRedis()
	suite.cache = NewCacheRedis(redisdb)
	suite.service = NewBillingService(db, cfg, make(chan bool, 1), nil, nil, nil, nil, nil, suite.cache)
	err = suite.service.Init()

	if err != nil {
		suite.FailNow("Billing service initialization failed", "%v", err)
	}
}

func (suite *ProjectTestSuite) TearDownTest() {
	if err := suite.service.db.Drop(); err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	suite.service.db.Close()
}

func (suite *ProjectTestSuite) TestProject_GetProjectById_Ok() {
	c, err := suite.service.project.GetById(suite.project.Id)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), c)
	assert.Equal(suite.T(), suite.project.Id, c.Id)
}

func (suite *ProjectTestSuite) TestProject_GetProjectById_NotFound() {
	_, err := suite.service.project.GetById(bson.NewObjectId().Hex())

	assert.Error(suite.T(), err)
	assert.Errorf(suite.T(), err, fmt.Sprintf(errorNotFound, pkg.CollectionProject))
}
