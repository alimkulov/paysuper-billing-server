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

type ProductTestSuite struct {
	suite.Suite
	mapper productMapper
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (suite *ProductTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *ProductTestSuite) Test_Product_NewProductMapper() {
	mapper := NewProductMapper()
	assert.IsType(suite.T(), &productMapper{}, mapper)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Ok() {
	original := &billingpb.Product{}
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

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Error_Id() {
	original := &billingpb.Product{
		Id:         "test",
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Error_MerchantId() {
	original := &billingpb.Product{
		MerchantId: "test",
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Error_ProjectId() {
	original := &billingpb.Product{
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Error_CreatedAt() {
	original := &billingpb.Product{
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Error_UpdatedAt() {
	original := &billingpb.Product{
		CreatedAt:  nil,
		UpdatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapObjectToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.Product{
		UpdatedAt:  nil,
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapMgoToObject_Ok() {
	original := &MgoProduct{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *ProductTestSuite) Test_Product_MapMgoToObject_Error_CreatedAt() {
	original := &MgoProduct{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *ProductTestSuite) Test_Product_MapMgoToObject_Error_UpdatedAt() {
	original := &MgoProduct{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
