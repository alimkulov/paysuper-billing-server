package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type pageReviewMapper struct{}

func NewPageReviewMapper() Mapper {
	return &pageReviewMapper{}
}

type MgoPageReview struct {
	Id        primitive.ObjectID `bson:"_id" faker:"objectId"`
	UserId    string             `bson:"user_id"`
	Review    string             `bson:"review"`
	Url       string             `bson:"url"`
	IsRead    bool               `bson:"is_read"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (m *pageReviewMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.PageReview)

	out := &MgoPageReview{
		UserId: in.UserId,
		Review: in.Review,
		Url:    in.Url,
		IsRead: in.IsRead,
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

func (m *pageReviewMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoPageReview)

	out := &billingpb.PageReview{
		Id:     in.Id.Hex(),
		UserId: in.UserId,
		Review: in.Review,
		Url:    in.Url,
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
