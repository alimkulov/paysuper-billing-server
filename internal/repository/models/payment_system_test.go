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

type PaymentSystemTestSuite struct {
	suite.Suite
	mapper paymentSystemMapper
}

func TestPaymentSystemTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentSystemTestSuite))
}

func (suite *PaymentSystemTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_NewPaymentSystemMapper() {
	mapper := NewPaymentSystemMapper()
	assert.IsType(suite.T(), &paymentSystemMapper{}, mapper)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapPaymentSystemToMgo_Ok() {
	original := &billingpb.PaymentSystem{}
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

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapPaymentSystemToMgo_Error_Id() {
	original := &billingpb.PaymentSystem{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapPaymentSystemToMgo_Error_CreatedAt() {
	original := &billingpb.PaymentSystem{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapPaymentSystemToMgo_Error_UpdatedAt() {
	original := &billingpb.PaymentSystem{
		CreatedAt: nil,
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapPaymentSystemToMgo_Ok_EmptyUpdatedAt() {
	original := &billingpb.PaymentSystem{
		UpdatedAt: nil,
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.NoError(suite.T(), err)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapMgoToPaymentSystem_Ok() {
	original := &MgoPaymentSystem{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapMgoToPaymentSystem_Error_CreatedAt() {
	original := &MgoPaymentSystem{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *PaymentSystemTestSuite) Test_PaymentSystem_MapMgoToPaymentSystem_Error_UpdatedAt() {
	original := &MgoPaymentSystem{
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
