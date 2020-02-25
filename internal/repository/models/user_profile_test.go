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

type UserProfileTestSuite struct {
	suite.Suite
	mapper userProfileMapper
}

func TestUserProfileTestSuite(t *testing.T) {
	suite.Run(t, new(UserProfileTestSuite))
}

func (suite *UserProfileTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *UserProfileTestSuite) Test_UserProfile_NewUserProfileMapper() {
	mapper := NewUserProfileMapper()
	assert.IsType(suite.T(), &userProfileMapper{}, mapper)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapUserProfileToMgo_Ok() {
	original := &billingpb.UserProfile{}
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

func (suite *UserProfileTestSuite) Test_UserProfile_MapUserProfileToMgo_Error_Id() {
	original := &billingpb.UserProfile{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapUserProfileToMgo_Error_CreatedAt() {
	original := &billingpb.UserProfile{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapUserProfileToMgo_Error_UpdatedAt() {
	original := &billingpb.UserProfile{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapUserProfileToMgo_Error_EmailConfirmedAt() {
	original := &billingpb.UserProfile{
		CreatedAt: nil,
		UpdatedAt: nil,
		Email: &billingpb.UserProfileEmail{
			ConfirmedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapMgoToUserProfile_Ok() {
	original := &MgoUserProfile{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapMgoToUserProfile_Error_CreatedAt() {
	original := &MgoUserProfile{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapMgoToUserProfile_Error_UpdatedAt() {
	original := &MgoUserProfile{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *UserProfileTestSuite) Test_UserProfile_MapMgoToUserProfile_Error_EmailConfirmedAt() {
	original := &MgoUserProfile{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: &MgoUserProfileEmail{
			ConfirmedAt: time.Time{}.AddDate(-10000, 0, 0),
		},
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
