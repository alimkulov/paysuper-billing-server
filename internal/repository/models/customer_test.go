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

type CustomerTestSuite struct {
	suite.Suite
	mapper customerMapper
}

func TestCustomerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerTestSuite))
}

func (suite *CustomerTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *CustomerTestSuite) Test_Customer_NewCustomerMapper() {
	mapper := NewCustomerMapper()
	assert.IsType(suite.T(), &customerMapper{}, mapper)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Ok() {
	original := &billingpb.Customer{}
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

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Error_Id() {
	original := &billingpb.Customer{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Error_CreatedAt() {
	original := &billingpb.Customer{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Error_UpdatedAt() {
	original := &billingpb.Customer{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.Customer{
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Error_IdentityMerchantId() {
	original := &billingpb.Customer{
		Identity: []*billingpb.CustomerIdentity{
			{MerchantId: "id"},
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapCustomerToMgo_Error_IdentityProjectId() {
	original := &billingpb.Customer{
		Identity: []*billingpb.CustomerIdentity{
			{MerchantId: primitive.NewObjectID().Hex(), ProjectId: "id"},
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapMgoToCustomer_Ok() {
	original := &MgoCustomer{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *CustomerTestSuite) Test_Customer_MapMgoToCustomer_Error_CreatedAt() {
	original := &MgoCustomer{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *CustomerTestSuite) Test_Customer_MapMgoToCustomer_Error_UpdatedAt() {
	original := &MgoCustomer{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
