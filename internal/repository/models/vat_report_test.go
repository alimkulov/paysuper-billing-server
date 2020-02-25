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

type VatReportTestSuite struct {
	suite.Suite
	mapper vatReportMapper
}

func TestVatReportTestSuite(t *testing.T) {
	suite.Run(t, new(VatReportTestSuite))
}

func (suite *VatReportTestSuite) SetupTest() {
	InitFakeCustomProviders()
}

func (suite *VatReportTestSuite) Test_VatReport_NewVatReportMapper() {
	mapper := NewVatReportMapper()
	assert.IsType(suite.T(), &vatReportMapper{}, mapper)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Ok() {
	original := &billingpb.VatReport{}
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

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_Id() {
	original := &billingpb.VatReport{
		Id: "test",
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_CreatedAt() {
	original := &billingpb.VatReport{
		CreatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_UpdatedAt() {
	original := &billingpb.VatReport{
		UpdatedAt: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_DateFrom() {
	original := &billingpb.VatReport{
		DateFrom: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_DateTo() {
	original := &billingpb.VatReport{
		DateFrom: &timestamp.Timestamp{},
		DateTo:   &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_PayUntilDate() {
	original := &billingpb.VatReport{
		DateFrom:     &timestamp.Timestamp{},
		DateTo:       &timestamp.Timestamp{},
		PayUntilDate: &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapVatReportToMgo_Error_PaidAt() {
	original := &billingpb.VatReport{
		DateFrom:     &timestamp.Timestamp{},
		DateTo:       &timestamp.Timestamp{},
		PayUntilDate: &timestamp.Timestamp{},
		PaidAt:       &timestamp.Timestamp{Seconds: -1, Nanos: -1},
	}
	_, err := suite.mapper.MapObjectToMgo(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Ok() {
	original := &MgoVatReport{}
	err := faker.FakeData(original)
	assert.NoError(suite.T(), err)

	obj, err := suite.mapper.MapMgoToObject(original)
	assert.NoError(suite.T(), err)

	mgo, err := suite.mapper.MapObjectToMgo(obj)
	assert.NoError(suite.T(), err)

	assert.ObjectsAreEqualValues(original, mgo)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_CreatedAt() {
	original := &MgoVatReport{
		CreatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_UpdatedAt() {
	original := &MgoVatReport{
		UpdatedAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_DateFrom() {
	original := &MgoVatReport{
		DateFrom: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_DateTo() {
	original := &MgoVatReport{
		DateTo: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_PayUntilDate() {
	original := &MgoVatReport{
		PayUntilDate: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}

func (suite *VatReportTestSuite) Test_VatReport_MapMgoToVatReport_Error_PaidAt() {
	original := &MgoVatReport{
		PaidAt: time.Time{}.AddDate(-10000, 0, 0),
	}
	_, err := suite.mapper.MapMgoToObject(original)
	assert.Error(suite.T(), err)
}
