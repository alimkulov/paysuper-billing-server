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

type ZipCodeTestSuite struct {
	suite.Suite
	mapper zipCodeMapper
}

func TestZipCodeTestSuite(t *testing.T) {
	suite.Run(t, new(ZipCodeTestSuite))
}

func (suite *ZipCodeTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *ZipCodeTestSuite) Test_UserProfile_NewZipCodeMapper() {
	mapper := NewZipCodeMapper()
	assert.IsType(suite.T(), &zipCodeMapper{}, mapper)
}

func (suite *ZipCodeTestSuite) Test_ZipCode_MapZipCodeToMgo_Ok() {
	original := &billingpb.ZipCode{}
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

func (suite *ZipCodeTestSuite) Test_ZipCode_MapZipCodeToMgo_Ok_EmptyCreatedAt() {
	original := &billingpb.ZipCode{}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *ZipCodeTestSuite) Test_ZipCode_MapZipCodeToMgo_Error_CreatedAt() {
	original := &billingpb.ZipCode{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *ZipCodeTestSuite) Test_ZipCode_MapMgoToZipCode_Ok() {
	original := &MgoZipCode{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *ZipCodeTestSuite) Test_ZipCode_MapMgoToZipCode_Error_CreatedAt() {
	original := &MgoZipCode{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
