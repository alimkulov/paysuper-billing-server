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

type MoneyBackCostMerchantTestSuite struct {
	suite.Suite
	mapper moneyBackCostMerchantMapper
}

func TestMoneyBackCostMerchantTestSuite(t *testing.T) {
	suite.Run(t, new(MoneyBackCostMerchantTestSuite))
}

func (suite *MoneyBackCostMerchantTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_NewMoneyBackCostMerchantMapper() {
	mapper := NewMoneyBackCostMerchantMapper()
	assert.IsType(suite.T(), &moneyBackCostMerchantMapper{}, mapper)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Ok() {
	original := &billingpb.MoneyBackCostMerchant{}
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

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Error_Id() {
	original := &billingpb.MoneyBackCostMerchant{
		Id:         "test",
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Error_MerchantId() {
	original := &billingpb.MoneyBackCostMerchant{
		MerchantId: "id",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Error_CreatedAt() {
	original := &billingpb.MoneyBackCostMerchant{
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Error_UpdatedAt() {
	original := &billingpb.MoneyBackCostMerchant{
		CreatedAt:  nil,
		UpdatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMoneyBackCostMerchantToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.MoneyBackCostMerchant{
		UpdatedAt:  nil,
		MerchantId: primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMgoToMoneyBackCostMerchant_Ok() {
	original := &MgoMoneyBackCostMerchant{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMgoToMoneyBackCostMerchant_Error_CreatedAt() {
	original := &MgoMoneyBackCostMerchant{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostMerchantTestSuite) Test_MoneyBackCostMerchant_MapMgoToMoneyBackCostMerchant_Error_UpdatedAt() {
	original := &MgoMoneyBackCostMerchant{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
