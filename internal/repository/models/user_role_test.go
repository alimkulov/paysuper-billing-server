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

type UserRoleTestSuite struct {
	suite.Suite
	mapper userRoleMapper
}

func TestUserRoleTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleTestSuite))
}

func (suite *UserRoleTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *UserRoleTestSuite) Test_UserRole_NewUserRoleMapper() {
	mapper := NewUserRoleMapper()
	assert.IsType(suite.T(), &userRoleMapper{}, mapper)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Ok() {
	original := &billingpb.UserRole{}
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

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Error_Id() {
	original := &billingpb.UserRole{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Error_CreatedAt() {
	original := &billingpb.UserRole{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Error_UpdatedAt() {
	original := &billingpb.UserRole{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Error_MerchantId() {
	original := &billingpb.UserRole{
		MerchantId: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapUserRoleToMgo_Error_UserId() {
	original := &billingpb.UserRole{
		UserId: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapMgoToUserRole_Ok() {
	original := &MgoUserRole{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapMgoToUserRole_Error_CreatedAt() {
	original := &MgoUserRole{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleTestSuite) Test_UserRole_MapMgoToUserRole_Error_UpdatedAt() {
	original := &MgoUserRole{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
