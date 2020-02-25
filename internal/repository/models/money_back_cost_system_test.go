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

type MoneyBackCostSystemTestSuite struct {
	suite.Suite
	mapper moneyBackCostSystemMapper
}

func TestMoneyBackCostSystemTestSuite(t *testing.T) {
	suite.Run(t, new(MoneyBackCostSystemTestSuite))
}

func (suite *MoneyBackCostSystemTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_NewMoneyBackCostSystemMapper() {
	mapper := NewMoneyBackCostSystemMapper()
	assert.IsType(suite.T(), &moneyBackCostSystemMapper{}, mapper)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMoneyBackCostSystemToMgo_Ok() {
	original := &billingpb.MoneyBackCostSystem{}
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

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMoneyBackCostSystemToMgo_Error_Id() {
	original := &billingpb.MoneyBackCostSystem{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMoneyBackCostSystemToMgo_Error_CreatedAt() {
	original := &billingpb.MoneyBackCostSystem{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMoneyBackCostSystemToMgo_Error_UpdatedAt() {
	original := &billingpb.MoneyBackCostSystem{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMoneyBackCostSystemToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.MoneyBackCostSystem{
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMgoToMoneyBackCostSystem_Ok() {
	original := &MgoMoneyBackCostSystem{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMgoToMoneyBackCostSystem_Error_CreatedAt() {
	original := &MgoMoneyBackCostSystem{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *MoneyBackCostSystemTestSuite) Test_MoneyBackCostSystem_MapMgoToMoneyBackCostSystem_Error_UpdatedAt() {
	original := &MgoMoneyBackCostSystem{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
