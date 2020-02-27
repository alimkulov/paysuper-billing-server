package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type userRoleMapper struct{}

func NewUserRoleMapper() Mapper {
	return &userRoleMapper{}
}

type MgoUserRole struct {
	Id         primitive.ObjectID  `bson:"_id" faker:"objectId"`
	MerchantId *primitive.ObjectID `bson:"merchant_id" faker:"objectIdPointer"`
	Role       string              `bson:"role"`
	Status     string              `bson:"status"`
	UserId     *primitive.ObjectID `bson:"user_id" faker:"objectIdPointer"`
	FirstName  string              `bson:"first_name"`
	LastName   string              `bson:"last_name"`
	Email      string              `bson:"email"`
	CreatedAt  time.Time           `bson:"created_at"`
	UpdatedAt  time.Time           `bson:"updated_at"`
}

func (m *userRoleMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.UserRole)

	out := &MgoUserRole{
		Role:      in.Role,
		Status:    in.Status,
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
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

	if in.MerchantId != "" {
		merchantOid, err := primitive.ObjectIDFromHex(in.MerchantId)

		if err != nil {
			return nil, err
		}

		out.MerchantId = &merchantOid
	}

	if in.UserId != "" {
		userOid, err := primitive.ObjectIDFromHex(in.UserId)

		if err != nil {
			return nil, err
		}

		out.UserId = &userOid
	}

	return out, nil
}

func (m *userRoleMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoUserRole)

	out := &billingpb.UserRole{
		Id:        in.Id.Hex(),
		Role:      in.Role,
		Status:    in.Status,
		Email:     in.Email,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if in.MerchantId != nil {
		out.MerchantId = in.MerchantId.Hex()
	}

	if in.UserId != nil {
		out.UserId = in.UserId.Hex()
	}

	return out, nil
}
