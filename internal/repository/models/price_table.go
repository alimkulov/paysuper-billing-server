package models

import (
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type priceTableMapper struct{}

func NewPriceTableMapper() Mapper {
	return &priceTableMapper{}
}

type MgoPriceTable struct {
	Id       primitive.ObjectID    `bson:"_id" faker:"objectId"`
	Currency string                `bson:"currency"`
	Ranges   []*MgoPriceTableRange `bson:"range"`
}

type MgoPriceTableRange struct {
	From     float64 `bson:"from"`
	To       float64 `bson:"to"`
	Position int32   `bson:"position"`
}

func (m *priceTableMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.PriceTable)

	out := &MgoPriceTable{
		Currency: in.Currency,
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

	if len(in.Ranges) > 0 {
		for _, v := range in.Ranges {
			out.Ranges = append(out.Ranges, &MgoPriceTableRange{From: v.From, To: v.To, Position: v.Position})
		}
	}

	return out, nil
}

func (m *priceTableMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	in := obj.(*MgoPriceTable)

	out := &billingpb.PriceTable{
		Id:       in.Id.Hex(),
		Currency: in.Currency,
	}

	rangesLen := len(in.Ranges)

	if rangesLen > 0 {
		out.Ranges = make([]*billingpb.PriceTableRange, rangesLen)

		for i, v := range in.Ranges {
			out.Ranges[i] = &billingpb.PriceTableRange{
				From:     v.From,
				To:       v.To,
				Position: v.Position,
			}
		}
	}

	return out, nil
}
