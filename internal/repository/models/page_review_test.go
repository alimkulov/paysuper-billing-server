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

type PageReviewTestSuite struct {
	suite.Suite
	mapper pageReviewMapper
}

func TestPageReviewTestSuite(t *testing.T) {
	suite.Run(t, new(PageReviewTestSuite))
}

func (suite *PageReviewTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PageReviewTestSuite) Test_PageReview_NewPageReviewMapper() {
	mapper := NewPageReviewMapper()
	assert.IsType(suite.T(), &pageReviewMapper{}, mapper)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapPageReviewToMgo_Ok() {
	original := &billingpb.PageReview{}
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

func (suite *PageReviewTestSuite) Test_PageReview_MapPageReviewToMgo_Error_Id() {
	original := &billingpb.PageReview{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapPageReviewToMgo_Error_CreatedAt() {
	original := &billingpb.PageReview{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapPageReviewToMgo_Error_UpdatedAt() {
	original := &billingpb.PageReview{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapPageReviewToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PageReview{
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapMgoToPageReview_Ok() {
	original := &MgoPageReview{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapMgoToPageReview_Error_CreatedAt() {
	original := &MgoPageReview{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PageReviewTestSuite) Test_PageReview_MapMgoToPageReview_Error_UpdatedAt() {
	original := &MgoPageReview{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
