package models

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MgoCountry struct {
	Id                      primitive.ObjectID             `bson:"_id" faker:"objectId"`
	IsoCodeA2               string                         `bson:"iso_code_a2"`
	Region                  string                         `bson:"region"`
	Currency                string                         `bson:"currency"`
	PaymentsAllowed         bool                           `bson:"payments_allowed"`
	ChangeAllowed           bool                           `bson:"change_allowed"`
	VatEnabled              bool                           `bson:"vat_enabled"`
	VatCurrency             string                         `bson:"vat_currency"`
	PriceGroupId            string                         `bson:"price_group_id"`
	VatThreshold            *billingpb.CountryVatThreshold `bson:"vat_threshold"`
	VatPeriodMonth          int32                          `bson:"vat_period_month"`
	VatDeadlineDays         int32                          `bson:"vat_deadline_days"`
	VatStoreYears           int32                          `bson:"vat_store_years"`
	VatCurrencyRatesPolicy  string                         `bson:"vat_currency_rates_policy"`
	VatCurrencyRatesSource  string                         `bson:"vat_currency_rates_source"`
	PayerTariffRegion       string                         `bson:"payer_tariff_region"`
	CreatedAt               time.Time                      `bson:"created_at"`
	UpdatedAt               time.Time                      `bson:"updated_at"`
	HighRiskPaymentsAllowed bool                           `bson:"high_risk_payments_allowed"`
	HighRiskChangeAllowed   bool                           `bson:"high_risk_change_allowed"`
}


type countryMapper struct{}

func (c *countryMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	m := obj.(*billingpb.Country)
	st := &MgoCountry{
		IsoCodeA2:               m.IsoCodeA2,
		Region:                  m.Region,
		Currency:                m.Currency,
		PaymentsAllowed:         m.PaymentsAllowed,
		ChangeAllowed:           m.ChangeAllowed,
		VatEnabled:              m.VatEnabled,
		PriceGroupId:            m.PriceGroupId,
		VatCurrency:             m.VatCurrency,
		VatThreshold:            m.VatThreshold,
		VatPeriodMonth:          m.VatPeriodMonth,
		VatStoreYears:           m.VatStoreYears,
		VatCurrencyRatesPolicy:  m.VatCurrencyRatesPolicy,
		VatCurrencyRatesSource:  m.VatCurrencyRatesSource,
		PayerTariffRegion:       m.PayerTariffRegion,
		HighRiskPaymentsAllowed: m.HighRiskPaymentsAllowed,
		HighRiskChangeAllowed:   m.HighRiskChangeAllowed,
	}
	if len(m.Id) <= 0 {
		st.Id = primitive.NewObjectID()
	} else {
		oid, err := primitive.ObjectIDFromHex(m.Id)

		if err != nil {
			return nil, errors.New(billingpb.ErrorInvalidObjectId)
		}

		st.Id = oid
	}

	if m.CreatedAt != nil {
		t, err := ptypes.Timestamp(m.CreatedAt)

		if err != nil {
			return nil, err
		}

		st.CreatedAt = t
	} else {
		st.CreatedAt = time.Now()
	}

	if m.UpdatedAt != nil {
		t, err := ptypes.Timestamp(m.UpdatedAt)

		if err != nil {
			return nil, err
		}

		st.UpdatedAt = t
	} else {
		st.UpdatedAt = time.Now()
	}

	return st, nil
}

func (c *countryMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	decoded := obj.(*MgoCountry)
	m := &billingpb.Country{}

	m.Id = decoded.Id.Hex()
	m.IsoCodeA2 = decoded.IsoCodeA2
	m.Region = decoded.Region
	m.Currency = decoded.Currency
	m.PaymentsAllowed = decoded.PaymentsAllowed
	m.ChangeAllowed = decoded.ChangeAllowed
	m.VatEnabled = decoded.VatEnabled
	m.PriceGroupId = decoded.PriceGroupId
	m.VatCurrency = decoded.VatCurrency
	m.VatThreshold = decoded.VatThreshold
	m.VatPeriodMonth = decoded.VatPeriodMonth
	m.VatDeadlineDays = decoded.VatDeadlineDays
	m.VatStoreYears = decoded.VatStoreYears
	m.VatCurrencyRatesPolicy = decoded.VatCurrencyRatesPolicy
	m.VatCurrencyRatesSource = decoded.VatCurrencyRatesSource
	m.PayerTariffRegion = decoded.PayerTariffRegion
	m.HighRiskPaymentsAllowed = decoded.HighRiskPaymentsAllowed
	m.HighRiskChangeAllowed = decoded.HighRiskChangeAllowed

	m.CreatedAt, err = ptypes.TimestampProto(decoded.CreatedAt)

	if err != nil {
		return nil, err
	}

	m.UpdatedAt, err = ptypes.TimestampProto(decoded.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func NewCountryMapper() Mapper {
	return &countryMapper{}
}
