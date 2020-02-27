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

type PaymentMinLimitSystemTestSuite struct {
	suite.Suite
	mapper paymentMinLimitSystemMapper
}

func TestPaymentMinLimitSystemTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentMinLimitSystemTestSuite))
}

func (suite *PaymentMinLimitSystemTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_NewPaymentMinLimitSystemMapper() {
	mapper := NewPaymentMinLimitSystemMapper()
	assert.IsType(suite.T(), &paymentMinLimitSystemMapper{}, mapper)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapPaymentMinLimitSystemToMgo_Ok() {
	original := &billingpb.PaymentMinLimitSystem{}
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

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapPaymentMinLimitSystemToMgo_Error_Id() {
	original := &billingpb.PaymentMinLimitSystem{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapPaymentMinLimitSystemToMgo_Error_CreatedAt() {
	original := &billingpb.PaymentMinLimitSystem{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapPaymentMinLimitSystemToMgo_Error_UpdatedAt() {
	original := &billingpb.PaymentMinLimitSystem{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapPaymentMinLimitSystemToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PaymentMinLimitSystem{
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapMgoToPaymentMinLimitSystem_Ok() {
	original := &MgoPaymentMinLimitSystem{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapMgoToPaymentMinLimitSystem_Error_CreatedAt() {
	original := &MgoPaymentMinLimitSystem{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentMinLimitSystemTestSuite) Test_PaymentMinLimitSystem_MapMgoToPaymentMinLimitSystem_Error_UpdatedAt() {
	original := &MgoPaymentMinLimitSystem{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
