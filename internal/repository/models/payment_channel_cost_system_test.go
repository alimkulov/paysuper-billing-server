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

type PaymentChannelCostSystemTestSuite struct {
	suite.Suite
	mapper paymentChannelCostSystemMapper
}

func TestPaymentChannelCostSystemTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentChannelCostSystemTestSuite))
}

func (suite *PaymentChannelCostSystemTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_NewPaymentChannelCostSystemMapper() {
	mapper := NewPaymentChannelCostSystemMapper()
	assert.IsType(suite.T(), &paymentChannelCostSystemMapper{}, mapper)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapPaymentChannelCostSystemToMgo_Ok() {
	original := &billingpb.PaymentChannelCostSystem{}
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

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapPaymentChannelCostSystemToMgo_Error_Id() {
	original := &billingpb.PaymentChannelCostSystem{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapPaymentChannelCostSystemToMgo_Error_CreatedAt() {
	original := &billingpb.PaymentChannelCostSystem{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapPaymentChannelCostSystemToMgo_Error_UpdatedAt() {
	original := &billingpb.PaymentChannelCostSystem{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapPaymentChannelCostSystemToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PaymentChannelCostSystem{
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapMgoToPaymentChannelCostSystem_Ok() {
	original := &MgoPaymentChannelCostSystem{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapMgoToPaymentChannelCostSystem_Error_CreatedAt() {
	original := &MgoPaymentChannelCostSystem{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostSystemTestSuite) Test_PaymentChannelCostSystem_MapMgoToPaymentChannelCostSystem_Error_UpdatedAt() {
	original := &MgoPaymentChannelCostSystem{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
