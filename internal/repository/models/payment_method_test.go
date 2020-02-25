package models

import (
	"github.com/bxcodec/faker"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type PaymentMethodTestSuite struct {
	suite.Suite
	mapper paymentMethodMapper
}

func TestPaymentMethodTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentMethodTestSuite))
}

func (suite *PaymentMethodTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PaymentMethodTestSuite) Test_UserProfile_NewPaymentMethodMapper() {
	mapper := NewPaymentMethodMapper()
	assert.IsType(suite.T(), &paymentMethodMapper{}, mapper)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Ok() {
	original := &billingpb.PaymentMethod{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), mgo)

	obj, err := suite.mapper.MapMgoToObject(mgo)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), obj)

	assert.ObjectsAreEqualValues(original, obj)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Error_Id() {
	original := &billingpb.PaymentMethod{
		Id: "id",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Ok_EmptyCreatedAt() {
	original := &billingpb.PaymentMethod{}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Error_CreatedAt() {
	original := &billingpb.PaymentMethod{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Error_UpdatedAt() {
	original := &billingpb.PaymentMethod{
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Error_PaymentSystemId() {
	original := &billingpb.PaymentMethod{
		PaymentSystemId: "id",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_Ok_EmptyStruct() {
	original := &billingpb.PaymentMethod{}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_SkipDuplicateTestSettings() {
	original := &billingpb.PaymentMethod{
		TestSettings: map[string]*billingpb.PaymentMethodParams{
			"TST1": {Currency: "RU", MccCode: "code", OperatingCompanyId: "company"},
			"TST2": {Currency: "RU", MccCode: "code", OperatingCompanyId: "company"},
		},
	}
	mgo, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), mgo)

	settings := mgo.(*MgoPaymentMethod).TestSettings
	assert.Len(suite.T(), settings, 1)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapPaymentMethodToMgo_SkipDuplicateProductionSettings() {
	original := &billingpb.PaymentMethod{
		ProductionSettings: map[string]*billingpb.PaymentMethodParams{
			"TST1": {Currency: "RU", MccCode: "code", OperatingCompanyId: "company"},
			"TST2": {Currency: "RU", MccCode: "code", OperatingCompanyId: "company"},
		},
	}
	mgo, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), mgo)

	settings := mgo.(*MgoPaymentMethod).ProductionSettings
	assert.Len(suite.T(), settings, 1)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapMgoToPaymentMethod_Ok() {
	original := &MgoPaymentMethod{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapMgoToPaymentMethod_Error_CreatedAt() {
	original := &MgoPaymentMethod{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMethodTestSuite) Test_PaymentMethod_MapMgoToPaymentMethod_Error_UpdatedAt() {
	original := &MgoPaymentMethod{
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
