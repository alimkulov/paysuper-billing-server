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

type PayoutTestSuite struct {
	suite.Suite
	mapper payoutMapper
}

type PayoutChangesTestSuite struct {
	suite.Suite
	mapper payoutChangesMapper
}

func TestPayoutTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutTestSuite))
}

func TestPayoutChangesTestSuite(t *testing.T) {
	suite.Run(t, new(PayoutChangesTestSuite))
}

func (suite *PayoutTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_NewPayoutDocumentsMapper() {
	mapper := NewPayoutMapper()
	assert.IsType(suite.T(), &payoutMapper{}, mapper)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok() {
	original := &billingpb.PayoutDocument{}
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

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_Id() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		Id:         "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_MerchantId() {
	original := &billingpb.PayoutDocument{
		MerchantId: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyCreatedAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		CreatedAt:  nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_CreatedAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_UpdatedAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		CreatedAt:  nil,
		UpdatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		UpdatedAt:  nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_ArrivalDate() {
	original := &billingpb.PayoutDocument{
		MerchantId:  primitive.NewObjectID().Hex(),
		ArrivalDate: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyArrivalDate() {
	original := &billingpb.PayoutDocument{
		MerchantId:  primitive.NewObjectID().Hex(),
		ArrivalDate: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_PeriodFrom() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PeriodFrom: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyPeriodFrom() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PeriodFrom: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_PeriodTo() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PeriodTo:   &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyPeriodTo() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PeriodTo:   nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Error_PaidAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PaidAt:     &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapPayoutDocumentsToMgo_Ok_EmptyPaidAt() {
	original := &billingpb.PayoutDocument{
		MerchantId: primitive.NewObjectID().Hex(),
		PaidAt:     nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Ok() {
	original := &MgoPayoutDocument{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_CreatedAt() {
	original := &MgoPayoutDocument{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_UpdatedAt() {
	original := &MgoPayoutDocument{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_ArrivalDate() {
	original := &MgoPayoutDocument{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ArrivalDate: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_PeriodFrom() {
	original := &MgoPayoutDocument{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ArrivalDate: time.Now(),
		PeriodFrom:  time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_PeriodTo() {
	original := &MgoPayoutDocument{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ArrivalDate: time.Now(),
		PeriodFrom:  time.Now(),
		PeriodTo:    time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutTestSuite) Test_PayoutDocuments_MapMgoToPayoutDocuments_Error_PaidAt() {
	original := &MgoPayoutDocument{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ArrivalDate: time.Now(),
		PeriodFrom:  time.Now(),
		PeriodTo:    time.Now(),
		PaidAt:      time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_NewPayoutChangesMapper() {
	mapper := NewPayoutChangesMapper()
	assert.IsType(suite.T(), &payoutChangesMapper{}, mapper)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapPayoutChangesToMgo_Ok() {
	original := &billingpb.PayoutDocumentChanges{}
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

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapPayoutChangesToMgo_Error_Id() {
	original := &billingpb.PayoutDocumentChanges{
		PayoutDocumentId: primitive.NewObjectID().Hex(),
		Id:               "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapPayoutChangesToMgo_Error_PayoutDocumentId() {
	original := &billingpb.PayoutDocumentChanges{
		PayoutDocumentId: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapPayoutChangesToMgo_Ok_EmptyCreatedAt() {
	original := &billingpb.PayoutDocumentChanges{
		PayoutDocumentId: primitive.NewObjectID().Hex(),
		CreatedAt:        nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapPayoutChangesToMgo_Error_CreatedAt() {
	original := &billingpb.PayoutDocumentChanges{
		PayoutDocumentId: primitive.NewObjectID().Hex(),
		CreatedAt:        &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapMgoToPayoutChanges_Ok() {
	original := &MgoPayoutDocumentChanges{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PayoutChangesTestSuite) Test_PayoutDocuments_MapMgoToPayoutChanges_Error_CreatedAt() {
	original := &MgoPayoutDocumentChanges{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
