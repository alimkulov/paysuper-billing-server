package models

import (
	"github.com/bxcodec/faker"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type PaymentChannelCostMerchantTestSuite struct {
	suite.Suite
	mapper paymentChannelCostMerchantMapper
}

func TestPaymentChannelCostMerchantTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentChannelCostMerchantTestSuite))
}

func (suite *PaymentChannelCostMerchantTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_NewPaymentChannelCostMerchantMapper() {
	mapper := NewPaymentChannelCostMerchantMapper()
	assert.IsType(suite.T(), &paymentChannelCostMerchantMapper{}, mapper)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Ok() {
	original := &billingpb.PaymentChannelCostMerchant{}
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

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Error_Id() {
	original := &billingpb.PaymentChannelCostMerchant{
		Id:         "test",
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Error_MerchantId() {
	original := &billingpb.PaymentChannelCostMerchant{
		MerchantId: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Error_CreatedAt() {
	original := &billingpb.PaymentChannelCostMerchant{
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Error_UpdatedAt() {
	original := &billingpb.PaymentChannelCostMerchant{
		CreatedAt:  nil,
		UpdatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapObjectToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PaymentChannelCostMerchant{
		UpdatedAt:  nil,
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapMgoToObject_Ok() {
	original := &MgoPaymentChannelCostMerchant{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapMgoToObject_Error_CreatedAt() {
	original := &MgoPaymentChannelCostMerchant{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentChannelCostMerchantTestSuite) Test_PaymentChannelCostMerchant_MapMgoToObject_Error_UpdatedAt() {
	original := &MgoPaymentChannelCostMerchant{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
