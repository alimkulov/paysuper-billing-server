package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type moneyBackCostSystemMapper struct{}

func NewMoneyBackCostSystemMapper() Mapper {
	return &moneyBackCostSystemMapper{}
}

type MgoMoneyBackCostSystem struct {
	Id                 primitive.ObjectID `bson:"_id" faker:"objectId"`
	Name               string             `bson:"name"`
	PayoutCurrency     string             `bson:"payout_currency"`
	UndoReason         string             `bson:"undo_reason"`
	Region             string             `bson:"region"`
	Country            string             `bson:"country"`
	DaysFrom           int32              `bson:"days_from"`
	PaymentStage       int32              `bson:"payment_stage"`
	Percent            float64            `bson:"percent"`
	FixAmount          float64            `bson:"fix_amount"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
	IsActive           bool               `bson:"is_active"`
	MccCode            string             `bson:"mcc_code"`
	OperatingCompanyId string             `bson:"operating_company_id"`
	FixAmountCurrency  string             `bson:"fix_amount_currency"`
}

func (m *moneyBackCostSystemMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.MoneyBackCostSystem)

	out := &MgoMoneyBackCostSystem{
		Name:               in.Name,
		PayoutCurrency:     in.PayoutCurrency,
		UndoReason:         in.UndoReason,
		Region:             in.Region,
		Country:            in.Country,
		DaysFrom:           in.DaysFrom,
		PaymentStage:       in.PaymentStage,
		Percent:            in.Percent,
		FixAmount:          in.FixAmount,
		IsActive:           in.IsActive,
		MccCode:            in.MccCode,
		OperatingCompanyId: in.OperatingCompanyId,
		FixAmountCurrency:  in.FixAmountCurrency,
	}

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

	if in.UpdatedAt != nil {
		t, err := ptypes.Timestamp(in.UpdatedAt)

		if err != nil {
			return nil, err
		}

		out.UpdatedAt = t
	} else {
		out.UpdatedAt = time.Now()
	}

	return out, nil
}

func (m *moneyBackCostSystemMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoMoneyBackCostSystem)

	out := &billingpb.MoneyBackCostSystem{
		Id:                 in.Id.Hex(),
		Name:               in.Name,
		PayoutCurrency:     in.PayoutCurrency,
		UndoReason:         in.UndoReason,
		Region:             in.Region,
		Country:            in.Country,
		DaysFrom:           in.DaysFrom,
		PaymentStage:       in.PaymentStage,
		Percent:            in.Percent,
		FixAmount:          in.FixAmount,
		IsActive:           in.IsActive,
		MccCode:            in.MccCode,
		OperatingCompanyId: in.OperatingCompanyId,
		FixAmountCurrency:  in.FixAmountCurrency,
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return out, nil
}
