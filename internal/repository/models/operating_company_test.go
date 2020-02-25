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

type OperatingCompanyTestSuite struct {
	suite.Suite
	mapper operatingCompanyMapper
}

func TestOperatingCompanyTestSuite(t *testing.T) {
	suite.Run(t, new(OperatingCompanyTestSuite))
}

func (suite *OperatingCompanyTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_NewOperatingCompanyMapper() {
	mapper := NewOperatingCompanyMapper()
	assert.IsType(suite.T(), &operatingCompanyMapper{}, mapper)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapOperatingCompanyToMgo_Ok() {
	original := &billingpb.OperatingCompany{}
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

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapOperatingCompanyToMgo_Error_Id() {
	original := &billingpb.OperatingCompany{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapOperatingCompanyToMgo_Error_CreatedAt() {
	original := &billingpb.OperatingCompany{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapOperatingCompanyToMgo_Error_UpdatedAt() {
	original := &billingpb.OperatingCompany{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapOperatingCompanyToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.OperatingCompany{
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapMgoToOperatingCompany_Ok() {
	original := &MgoOperatingCompany{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapMgoToOperatingCompany_Error_CreatedAt() {
	original := &MgoOperatingCompany{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *OperatingCompanyTestSuite) Test_OperatingCompany_MapMgoToOperatingCompany_Error_UpdatedAt() {
	original := &MgoOperatingCompany{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
