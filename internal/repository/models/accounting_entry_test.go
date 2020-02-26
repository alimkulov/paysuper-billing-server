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

type AccountingEntryTestSuite struct {
	suite.Suite
	mapper accountingEntryMapper
}

func TestAccountingEntryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountingEntryTestSuite))
}

func (suite *AccountingEntryTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_NewAccountingEntryMapper() {
	mapper := NewAccountingEntryMapper()
	assert.IsType(suite.T(), &accountingEntryMapper{}, mapper)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Ok() {
	original := &billingpb.AccountingEntry{}
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

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Error_Id() {
	original := &billingpb.AccountingEntry{
		Id:         "test",
		MerchantId: primitive.NewObjectID().Hex(),
		Source: &billingpb.AccountingEntrySource{
			Id:   primitive.NewObjectID().Hex(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Error_MerchantId() {
	original := &billingpb.AccountingEntry{
		MerchantId: "test",
		Source: &billingpb.AccountingEntrySource{
			Id:   primitive.NewObjectID().Hex(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Error_SourceId() {
	original := &billingpb.AccountingEntry{
		MerchantId: primitive.NewObjectID().Hex(),
		Source: &billingpb.AccountingEntrySource{
			Id:   "id",
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Error_CreatedAt() {
	original := &billingpb.AccountingEntry{
		CreatedAt:  &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId: primitive.NewObjectID().Hex(),
		Source: &billingpb.AccountingEntrySource{
			Id:   primitive.NewObjectID().Hex(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Error_AvailableOn() {
	original := &billingpb.AccountingEntry{
		CreatedAt:   nil,
		AvailableOn: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
		MerchantId:  primitive.NewObjectID().Hex(),
		Source: &billingpb.AccountingEntrySource{
			Id:   primitive.NewObjectID().Hex(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapObjectToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.AccountingEntry{
		AvailableOn: nil,
		MerchantId:  primitive.NewObjectID().Hex(),
		Source: &billingpb.AccountingEntrySource{
			Id:   primitive.NewObjectID().Hex(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapMgoToObject_Ok() {
	original := &MgoAccountingEntry{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapMgoToObject_Error_CreatedAt() {
	original := &MgoAccountingEntry{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
		Source: &MgoAccountingEntrySource{
			Id:   primitive.NewObjectID(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *AccountingEntryTestSuite) Test_AccountingEntry_MapMgoToObject_Error_AvailableOn() {
	original := &MgoAccountingEntry{
		CreatedAt:   time.Now(),
		AvailableOn: time.Time{}.AddDate(-10000, 0, 0),
		Source: &MgoAccountingEntrySource{
			Id:   primitive.NewObjectID(),
			Type: "type",
		},
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
