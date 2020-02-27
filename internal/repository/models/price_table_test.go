package models

import (
	"github.com/bxcodec/faker"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PriceTableTestSuite struct {
	suite.Suite
	mapper priceTableMapper
}

func TestPriceTableTestSuite(t *testing.T) {
	suite.Run(t, new(PriceTableTestSuite))
}

func (suite *PriceTableTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PriceTableTestSuite) Test_PriceTable_NewPriceTableMapper() {
	mapper := NewPriceTableMapper()
	assert.IsType(suite.T(), &priceTableMapper{}, mapper)
}

func (suite *PriceTableTestSuite) Test_PriceTable_MapObjectToMgo_Ok() {
	original := &billingpb.PriceTable{}
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

func (suite *PriceTableTestSuite) Test_PriceTable_MapObjectToMgo_Ok_EmptyId() {
	original := &billingpb.PriceTable{}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PriceTableTestSuite) Test_PriceTable_MapObjectToMgo_Error_Id() {
	original := &billingpb.PriceTable{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PriceTableTestSuite) Test_PriceTable_MapMgoToObject_Ok() {
	original := &MgoPriceTable{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}
