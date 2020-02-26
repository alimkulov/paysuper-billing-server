package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type accountingEntryMapper struct{}

func NewAccountingEntryMapper() Mapper {
	return &accountingEntryMapper{}
}

type MgoAccountingEntrySource struct {
	Id   primitive.ObjectID `bson:"id" faker:"objectId"`
	Type string             `bson:"type"`
}

type MgoAccountingEntry struct {
	Id                 primitive.ObjectID        `bson:"_id" faker:"objectId"`
	Object             string                    `bson:"object"`
	Type               string                    `bson:"type"`
	Source             *MgoAccountingEntrySource `bson:"source"`
	MerchantId         primitive.ObjectID        `bson:"merchant_id" faker:"objectId"`
	Amount             float64                   `bson:"amount"`
	Currency           string                    `bson:"currency"`
	Reason             string                    `bson:"reason"`
	Status             string                    `bson:"status"`
	Country            string                    `bson:"country"`
	OriginalAmount     float64                   `bson:"original_amount"`
	OriginalCurrency   string                    `bson:"original_currency"`
	LocalAmount        float64                   `bson:"local_amount"`
	LocalCurrency      string                    `bson:"local_currency"`
	CreatedAt          time.Time                 `bson:"created_at"`
	AvailableOn        time.Time                 `bson:"available_on"`
	OperatingCompanyId string                    `bson:"operating_company_id"`
}

func (m *accountingEntryMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.AccountingEntry)

	out := &MgoAccountingEntry{
		Object: in.Object,
		Type:   in.Type,
		Source: &MgoAccountingEntrySource{
			Type: in.Source.Type,
		},
		Amount:             in.Amount,
		Currency:           in.Currency,
		OriginalAmount:     in.OriginalAmount,
		OriginalCurrency:   in.OriginalCurrency,
		LocalAmount:        in.LocalAmount,
		LocalCurrency:      in.LocalCurrency,
		Country:            in.Country,
		Reason:             in.Reason,
		Status:             in.Status,
		OperatingCompanyId: in.OperatingCompanyId,
	}

	merchantOid, err := primitive.ObjectIDFromHex(in.MerchantId)

	if err != nil {
		return nil, err
	}

	out.MerchantId = merchantOid

	sourceOid, err := primitive.ObjectIDFromHex(in.Source.Id)

	if err != nil {
		return nil, err
	}

	out.Source.Id = sourceOid

	if len(in.Id) <= 0 {
		out.Id = primitive.NewObjectID()
	} else {
		oid, err := primitive.ObjectIDFromHex(in.Id)

		if err != nil {
			return nil, err
		}

		out.Id = oid
	}

	if in.CreatedAt != nil {
		t, err := ptypes.Timestamp(in.CreatedAt)

		if err != nil {
			return nil, err
		}

		out.CreatedAt = t
	} else {
		out.CreatedAt = time.Now()
	}

	if in.AvailableOn != nil {
		t, err := ptypes.Timestamp(in.AvailableOn)

		if err != nil {
			return nil, err
		}

		out.AvailableOn = t
	}

	return out, nil
}

func (m *accountingEntryMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoAccountingEntry)

	out := &billingpb.AccountingEntry{
		Id:     in.Id.Hex(),
		Object: in.Object,
		Type:   in.Type,
		Source: &billingpb.AccountingEntrySource{
			Id:   in.Source.Id.Hex(),
			Type: in.Source.Type,
		},
		MerchantId:         in.MerchantId.Hex(),
		Amount:             in.Amount,
		Currency:           in.Currency,
		OriginalAmount:     in.OriginalAmount,
		OriginalCurrency:   in.OriginalCurrency,
		LocalAmount:        in.LocalAmount,
		LocalCurrency:      in.LocalCurrency,
		Country:            in.Country,
		Reason:             in.Reason,
		Status:             in.Status,
		OperatingCompanyId: in.OperatingCompanyId,
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.AvailableOn, err = ptypes.TimestampProto(in.AvailableOn)

	if err != nil {
		return nil, err
	}

	return out, nil
}
