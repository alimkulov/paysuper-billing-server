package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type operatingCompanyMapper struct{}

func NewOperatingCompanyMapper() Mapper {
	return &operatingCompanyMapper{}
}

type MgoOperatingCompany struct {
	Id                 primitive.ObjectID `bson:"_id" faker:"objectId"`
	Name               string             `bson:"name"`
	Country            string             `bson:"country"`
	RegistrationNumber string             `bson:"registration_number"`
	RegistrationDate   string             `bson:"registration_date"`
	VatNumber          string             `bson:"vat_number"`
	Email              string             `bson:"email"`
	Address            string             `bson:"address"`
	VatAddress         string             `bson:"vat_address"`
	SignatoryName      string             `bson:"signatory_name"`
	SignatoryPosition  string             `bson:"signatory_position"`
	BankingDetails     string             `bson:"banking_details"`
	PaymentCountries   []string           `bson:"payment_countries"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

func (m *operatingCompanyMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.OperatingCompany)

	out := &MgoOperatingCompany{
		Name:               in.Name,
		Country:            in.Country,
		RegistrationNumber: in.RegistrationNumber,
		RegistrationDate:   in.RegistrationDate,
		VatNumber:          in.VatNumber,
		Email:              in.Email,
		Address:            in.Address,
		VatAddress:         in.VatAddress,
		SignatoryName:      in.SignatoryName,
		SignatoryPosition:  in.SignatoryPosition,
		BankingDetails:     in.BankingDetails,
		PaymentCountries:   in.PaymentCountries,
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

func (m *operatingCompanyMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoOperatingCompany)

	out := &billingpb.OperatingCompany{
		Id:                 in.Id.Hex(),
		Name:               in.Name,
		Country:            in.Country,
		RegistrationNumber: in.RegistrationNumber,
		RegistrationDate:   in.RegistrationDate,
		VatNumber:          in.VatNumber,
		Email:              in.Email,
		Address:            in.Address,
		VatAddress:         in.VatAddress,
		SignatoryName:      in.SignatoryName,
		SignatoryPosition:  in.SignatoryPosition,
		BankingDetails:     in.BankingDetails,
		PaymentCountries:   in.PaymentCountries,
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
