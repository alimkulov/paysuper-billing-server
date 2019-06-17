package service

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-billing-server/internal/config"
	"github.com/paysuper/paysuper-billing-server/internal/mock"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	mongodb "github.com/paysuper/paysuper-database-mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ZipCodeTestSuite struct {
	suite.Suite
	service *Service
}

func Test_ZipCode(t *testing.T) {
	suite.Run(t, new(ZipCodeTestSuite))
}

func (suite *ZipCodeTestSuite) SetupTest() {
	cfg, err := config.NewConfig()

	if err != nil {
		suite.FailNow("Config load failed", "%v", err)
	}

	cfg.AccountingCurrency = "RUB"
	db, err := mongodb.NewDatabase()

	if err != nil {
		suite.FailNow("Database connection failed", "%v", err)
	}

	zipCode := &billing.ZipCode{
		Zip:     "98001",
		Country: "US",
		City:    "Washington",
		State: &billing.ZipCodeState{
			Code: "NJ",
			Name: "New Jersey",
		},
		CreatedAt: ptypes.TimestampNow(),
	}

	err = db.Collection(collectionZipCode).Insert(zipCode)

	if err != nil {
		suite.FailNow("Insert zip codes test data failed", "%v", err)
	}

	rub := &billing.Currency{
		CodeInt:  643,
		CodeA3:   "RUB",
		Name:     &billing.Name{Ru: "Российский рубль", En: "Russian ruble"},
		IsActive: true,
	}
	err = InitTestCurrency(db, []interface{}{rub})

	if err != nil {
		suite.FailNow("Insert currency test data failed", "%v", err)
	}

	redisdb := mock.NewTestRedis()
	cache := NewCacheRedis(redisdb)
	suite.service = NewBillingService(db, cfg, nil, nil, nil, nil, nil, cache)
	err = suite.service.Init()

	if err != nil {
		suite.FailNow("Billing service initialization failed", "%v", err)
	}
}

func (suite *ZipCodeTestSuite) TearDownTest() {
	err := suite.service.db.Drop()

	if err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	suite.service.db.Close()
}

func (suite *ZipCodeTestSuite) TestZipCode_GetExist_Ok() {
	zip := "98001"
	zipCode, err := suite.service.zipCode.getByZipAndCountry(zip, CountryCodeUSA)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), zipCode)
	assert.Equal(suite.T(), zip, zipCode.Zip)
	assert.Equal(suite.T(), CountryCodeUSA, zipCode.Country)

	zipCode, err = suite.service.zipCode.getByZipAndCountry(zip, CountryCodeUSA)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), zipCode)
	assert.Equal(suite.T(), zip, zipCode.Zip)
	assert.Equal(suite.T(), CountryCodeUSA, zipCode.Country)
}

func (suite *ZipCodeTestSuite) TestZipCode_NotFound_Error() {
	zip := "98002"
	zipCode, err := suite.service.zipCode.getByZipAndCountry(zip, CountryCodeUSA)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), fmt.Sprintf(errorNotFound, collectionZipCode), err.Error())
	assert.Nil(suite.T(), zipCode)
}
