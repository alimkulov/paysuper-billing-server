package models

import (
	"github.com/bxcodec/faker"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NotificationTestSuite struct {
	suite.Suite
	mapper notificationMapper
}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (suite *NotificationTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *NotificationTestSuite) Test_UserProfile_NewZipCodeMapper() {
	mapper := NewNotificationMapper()
	assert.IsType(suite.T(), &notificationMapper{}, mapper)
}

func (suite *NotificationTestSuite) Test_Notification_MapToMgo_Ok() {
	original := &billingpb.Notification{}
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

func (suite *NotificationTestSuite) Test_Notification_Error_CreatedAt() {
	original := &billingpb.Notification{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}
