package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type paymentMinLimitSystemMapper struct{}

func NewPaymentMinLimitSystemMapper() Mapper {
	return &paymentMinLimitSystemMapper{}
}

type MgoPaymentMinLimitSystem struct {
	Id        primitive.ObjectID `bson:"_id" faker:"objectId"`
	Currency  string             `bson:"currency"`
	Amount    float64            `bson:"amount"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (m *paymentMinLimitSystemMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.PaymentMinLimitSystem)

	out := &MgoPaymentMinLimitSystem{
		Currency: in.Currency,
		Amount:   in.Amount,
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

func (m *paymentMinLimitSystemMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoPaymentMinLimitSystem)

	out := &billingpb.PaymentMinLimitSystem{
		Id:       in.Id.Hex(),
		Currency: in.Currency,
		Amount:   in.Amount,
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
