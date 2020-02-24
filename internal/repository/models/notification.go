package models

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MgoNotification struct {
	Id         primitive.ObjectID                    `bson:"_id"`
	Message    string                                `bson:"message"`
	MerchantId primitive.ObjectID                    `bson:"merchant_id"`
	UserId     string                                `bson:"user_id"`
	IsSystem   bool                                  `bson:"is_system"`
	IsRead     bool                                  `bson:"is_read"`
	CreatedAt  time.Time                             `bson:"created_at"`
	UpdatedAt  time.Time                             `bson:"updated_at"`
	Statuses   *billingpb.SystemNotificationStatuses `bson:"statuses"`
}

type notificationMapper struct {

}

func (n *notificationMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	m := obj.(*billingpb.Notification)
	oid, err := primitive.ObjectIDFromHex(m.MerchantId)

	if err != nil {
		return nil, err
	}

	st := &MgoNotification{
		Message:    m.Message,
		IsSystem:   m.IsSystem,
		IsRead:     m.IsRead,
		MerchantId: oid,
		UserId:     m.UserId,
		Statuses:   m.Statuses,
	}

	if len(m.Id) <= 0 {
		st.Id = primitive.NewObjectID()
	} else {
		oid, err = primitive.ObjectIDFromHex(m.Id)

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
		m.CreatedAt, _ = ptypes.TimestampProto(st.CreatedAt)
	}

	if m.UpdatedAt != nil {
		t, err := ptypes.Timestamp(m.UpdatedAt)

		if err != nil {
			return nil, err
		}

		st.UpdatedAt = t
	} else {
		st.UpdatedAt = time.Now()
		m.UpdatedAt, _ = ptypes.TimestampProto(st.UpdatedAt)
	}

	return st, nil
}

func (n *notificationMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	decoded := obj.(*MgoNotification)
	m := &billingpb.Notification{}

	m.Id = decoded.Id.Hex()
	m.Message = decoded.Message
	m.IsSystem = decoded.IsSystem
	m.IsRead = decoded.IsRead
	m.MerchantId = decoded.MerchantId.Hex()
	m.UserId = decoded.UserId
	m.Statuses = decoded.Statuses

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

func NewNotificationMapper() Mapper {
	return &notificationMapper{}
}
