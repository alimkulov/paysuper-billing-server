package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type customerMapper struct{}

func NewCustomerMapper() Mapper {
	return &customerMapper{}
}

type MgoCustomer struct {
	Id                    primitive.ObjectID               `bson:"_id" faker:"objectId"`
	TechEmail             string                           `bson:"tech_email"`
	ExternalId            string                           `bson:"external_id"`
	Email                 string                           `bson:"email"`
	EmailVerified         bool                             `bson:"email_verified"`
	Phone                 string                           `bson:"phone"`
	PhoneVerified         bool                             `bson:"phone_verified"`
	Name                  string                           `bson:"name"`
	Ip                    []byte                           `bson:"ip"`
	Locale                string                           `bson:"locale"`
	AcceptLanguage        string                           `bson:"accept_language"`
	UserAgent             string                           `bson:"user_agent"`
	Address               *billingpb.OrderBillingAddress   `bson:"address"`
	Identity              []*MgoCustomerIdentity           `bson:"identity"`
	IpHistory             []*MgoCustomerIpHistory          `bson:"ip_history"`
	AddressHistory        []*MgoCustomerAddressHistory     `bson:"address_history"`
	LocaleHistory         []*MgoCustomerStringValueHistory `bson:"locale_history"`
	AcceptLanguageHistory []*MgoCustomerStringValueHistory `bson:"accept_language_history"`
	Metadata              map[string]string                `bson:"metadata"`
	CreatedAt             time.Time                        `bson:"created_at"`
	UpdatedAt             time.Time                        `bson:"updated_at"`
	NotifySale            bool                             `bson:"notify_sale"`
	NotifySaleEmail       string                           `bson:"notify_sale_email"`
	NotifyNewRegion       bool                             `bson:"notify_new_region"`
	NotifyNewRegionEmail  string                           `bson:"notify_new_region_email"`
}

type MgoCustomerIdentity struct {
	MerchantId primitive.ObjectID `bson:"merchant_id" faker:"objectId"`
	ProjectId  primitive.ObjectID `bson:"project_id" faker:"objectId"`
	Type       string             `bson:"type"`
	Value      string             `bson:"value"`
	Verified   bool               `bson:"verified"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type MgoCustomerIpHistory struct {
	Ip        []byte    `bson:"ip"`
	CreatedAt time.Time `bson:"created_at"`
}

type MgoCustomerAddressHistory struct {
	Country    string    `bson:"country"`
	City       string    `bson:"city"`
	PostalCode string    `bson:"postal_code"`
	State      string    `bson:"state"`
	CreatedAt  time.Time `bson:"created_at"`
}

type MgoCustomerStringValueHistory struct {
	Value     string    `bson:"value"`
	CreatedAt time.Time `bson:"created_at"`
}

func (m *customerMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.Customer)

	out := &MgoCustomer{
		TechEmail:             in.TechEmail,
		ExternalId:            in.ExternalId,
		Email:                 in.Email,
		EmailVerified:         in.EmailVerified,
		Phone:                 in.Phone,
		PhoneVerified:         in.PhoneVerified,
		Name:                  in.Name,
		Ip:                    in.Ip,
		Locale:                in.Locale,
		AcceptLanguage:        in.AcceptLanguage,
		UserAgent:             in.UserAgent,
		Address:               in.Address,
		Metadata:              in.Metadata,
		NotifySale:            in.NotifySale,
		NotifySaleEmail:       in.NotifySaleEmail,
		NotifyNewRegion:       in.NotifyNewRegion,
		NotifyNewRegionEmail:  in.NotifyNewRegionEmail,
		Identity:              []*MgoCustomerIdentity{},
		IpHistory:             []*MgoCustomerIpHistory{},
		AddressHistory:        []*MgoCustomerAddressHistory{},
		LocaleHistory:         []*MgoCustomerStringValueHistory{},
		AcceptLanguageHistory: []*MgoCustomerStringValueHistory{},
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

	for _, v := range in.Identity {
		merchantOid, err := primitive.ObjectIDFromHex(v.MerchantId)

		if err != nil {
			return nil, err
		}

		projectOid, err := primitive.ObjectIDFromHex(v.ProjectId)

		if err != nil {
			return nil, err
		}

		mgoIdentity := &MgoCustomerIdentity{
			MerchantId: merchantOid,
			ProjectId:  projectOid,
			Type:       v.Type,
			Value:      v.Value,
			Verified:   v.Verified,
		}

		mgoIdentity.CreatedAt, _ = ptypes.Timestamp(v.CreatedAt)
		out.Identity = append(out.Identity, mgoIdentity)
	}

	for _, v := range in.IpHistory {
		mgoIdentity := &MgoCustomerIpHistory{Ip: v.Ip}
		mgoIdentity.CreatedAt, _ = ptypes.Timestamp(v.CreatedAt)
		out.IpHistory = append(out.IpHistory, mgoIdentity)
	}

	for _, v := range in.AddressHistory {
		mgoIdentity := &MgoCustomerAddressHistory{
			Country:    v.Country,
			City:       v.City,
			PostalCode: v.PostalCode,
			State:      v.State,
		}
		mgoIdentity.CreatedAt, _ = ptypes.Timestamp(v.CreatedAt)
		out.AddressHistory = append(out.AddressHistory, mgoIdentity)
	}

	for _, v := range in.LocaleHistory {
		mgoIdentity := &MgoCustomerStringValueHistory{Value: v.Value}
		mgoIdentity.CreatedAt, _ = ptypes.Timestamp(v.CreatedAt)
		out.LocaleHistory = append(out.LocaleHistory, mgoIdentity)
	}

	for _, v := range in.AcceptLanguageHistory {
		mgoIdentity := &MgoCustomerStringValueHistory{Value: v.Value}
		mgoIdentity.CreatedAt, _ = ptypes.Timestamp(v.CreatedAt)
		out.AcceptLanguageHistory = append(out.AcceptLanguageHistory, mgoIdentity)
	}

	return out, nil
}

func (m *customerMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoCustomer)

	out := &billingpb.Customer{
		Id:                    in.Id.Hex(),
		TechEmail:             in.TechEmail,
		ExternalId:            in.ExternalId,
		Email:                 in.Email,
		EmailVerified:         in.EmailVerified,
		Phone:                 in.Phone,
		PhoneVerified:         in.PhoneVerified,
		Name:                  in.Name,
		Ip:                    in.Ip,
		Locale:                in.Locale,
		AcceptLanguage:        in.AcceptLanguage,
		UserAgent:             in.UserAgent,
		Address:               in.Address,
		Metadata:              in.Metadata,
		NotifySale:            in.NotifySale,
		NotifySaleEmail:       in.NotifySaleEmail,
		NotifyNewRegion:       in.NotifyNewRegion,
		NotifyNewRegionEmail:  in.NotifyNewRegionEmail,
		Identity:              []*billingpb.CustomerIdentity{},
		IpHistory:             []*billingpb.CustomerIpHistory{},
		AddressHistory:        []*billingpb.CustomerAddressHistory{},
		LocaleHistory:         []*billingpb.CustomerStringValueHistory{},
		AcceptLanguageHistory: []*billingpb.CustomerStringValueHistory{},
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt)

	if err != nil {
		return nil, err
	}

	for _, v := range in.Identity {
		identity := &billingpb.CustomerIdentity{
			MerchantId: v.MerchantId.Hex(),
			ProjectId:  v.ProjectId.Hex(),
			Type:       v.Type,
			Value:      v.Value,
			Verified:   v.Verified,
		}

		identity.CreatedAt, _ = ptypes.TimestampProto(v.CreatedAt)
		out.Identity = append(out.Identity, identity)
	}

	for _, v := range in.IpHistory {
		identity := &billingpb.CustomerIpHistory{Ip: v.Ip}
		identity.CreatedAt, _ = ptypes.TimestampProto(v.CreatedAt)
		out.IpHistory = append(out.IpHistory, identity)
	}

	for _, v := range in.AddressHistory {
		identity := &billingpb.CustomerAddressHistory{
			Country:    v.Country,
			City:       v.City,
			PostalCode: v.PostalCode,
			State:      v.State,
		}
		identity.CreatedAt, _ = ptypes.TimestampProto(v.CreatedAt)
		out.AddressHistory = append(out.AddressHistory, identity)
	}

	for _, v := range in.LocaleHistory {
		identity := &billingpb.CustomerStringValueHistory{Value: v.Value}
		identity.CreatedAt, _ = ptypes.TimestampProto(v.CreatedAt)
		out.LocaleHistory = append(out.LocaleHistory, identity)
	}

	for _, v := range in.AcceptLanguageHistory {
		identity := &billingpb.CustomerStringValueHistory{Value: v.Value}
		identity.CreatedAt, _ = ptypes.TimestampProto(v.CreatedAt)
		out.AcceptLanguageHistory = append(out.AcceptLanguageHistory, identity)
	}

	return out, nil
}
