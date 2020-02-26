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

type KeyProductTestSuite struct {
	suite.Suite
	mapper keyProductMapper
}

func TestKeyProductTestSuite(t *testing.T) {
	suite.Run(t, new(KeyProductTestSuite))
}

func (suite *KeyProductTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *KeyProductTestSuite) Test_KeyProduct_NewKeyProductMapper() {
	mapper := NewKeyProductMapper()
	assert.IsType(suite.T(), &keyProductMapper{}, mapper)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Ok() {
	original := &billingpb.KeyProduct{}
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

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_Id() {
	original := &billingpb.KeyProduct{
		Id:         "test",
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_MerchantId() {
	original := &billingpb.KeyProduct{
		MerchantId: "test",
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_ProjectId() {
	original := &billingpb.KeyProduct{
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_CreatedAt() {
	original := &billingpb.KeyProduct{
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_UpdatedAt() {
	original := &billingpb.KeyProduct{
		CreatedAt:  nil,
		UpdatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.KeyProduct{
		UpdatedAt:  nil,
		MerchantId: primitive.NewObjectID().Hex(),
		ProjectId:  primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapObjectToMgo_Error_PublishedAt() {
	original := &billingpb.KeyProduct{
		CreatedAt:   nil,
		UpdatedAt:   nil,
		PublishedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId:  primitive.NewObjectID().Hex(),
		ProjectId:   primitive.NewObjectID().Hex(),
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapMgoToObject_Ok() {
	original := &MgoKeyProduct{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapMgoToObject_Error_CreatedAt() {
	original := &MgoKeyProduct{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapMgoToObject_Error_UpdatedAt() {
	original := &MgoKeyProduct{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *KeyProductTestSuite) Test_KeyProduct_MapMgoToObject_Error_PublishedAt() {
	t := time.Time{}.AddDate(-10000, 0, 0)
	original := &MgoKeyProduct{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		PublishedAt: &t,
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
