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

type PriceGroupTestSuite struct {
	suite.Suite
	mapper priceGroupMapper
}

func TestPriceGroupTestSuite(t *testing.T) {
	suite.Run(t, new(PriceGroupTestSuite))
}

func (suite *PriceGroupTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_NewPriceGroupMapper() {
	mapper := NewPriceGroupMapper()
	assert.IsType(suite.T(), &priceGroupMapper{}, mapper)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapPriceGroupToMgo_Ok() {
	original := &billingpb.PriceGroup{}
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

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapPriceGroupToMgo_Error_Id() {
	original := &billingpb.PriceGroup{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapPriceGroupToMgo_Error_CreatedAt() {
	original := &billingpb.PriceGroup{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapPriceGroupToMgo_Error_UpdatedAt() {
	original := &billingpb.PriceGroup{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapPriceGroupToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PriceGroup{
		CreatedAt: nil,
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapMgoToPriceGroup_Ok() {
	original := &MgoPriceGroup{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapMgoToPriceGroup_Error_CreatedAt() {
	original := &MgoPriceGroup{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PriceGroupTestSuite) Test_PriceGroup_MapMgoToPriceGroup_Error_UpdatedAt() {
	original := &MgoPriceGroup{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
