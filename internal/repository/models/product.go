package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	tools "github.com/paysuper/paysuper-tools/number"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type productMapper struct{}

func NewProductMapper() Mapper {
	return &productMapper{}
}

type MgoProduct struct {
	Id              primitive.ObjectID              `bson:"_id" json:"id" faker:"objectId"`
	Object          string                          `bson:"object" json:"object"`
	Type            string                          `bson:"type" json:"type"`
	Sku             string                          `bson:"sku" json:"sku"`
	Name            []*billingpb.I18NTextSearchable `bson:"name" json:"name"`
	DefaultCurrency string                          `bson:"default_currency" json:"default_currency"`
	Enabled         bool                            `bson:"enabled" json:"enabled"`
	Prices          []*billingpb.ProductPrice       `bson:"prices" json:"prices"`
	Description     map[string]string               `bson:"description" json:"description"`
	LongDescription map[string]string               `bson:"long_description,omitempty" json:"long_description"`
	CreatedAt       time.Time                       `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time                       `bson:"updated_at" json:"updated_at"`
	Images          []string                        `bson:"images,omitempty" json:"images"`
	Url             string                          `bson:"url,omitempty" json:"url"`
	Metadata        map[string]string               `bson:"metadata,omitempty" json:"metadata"`
	Deleted         bool                            `bson:"deleted" json:"deleted"`
	MerchantId      primitive.ObjectID              `bson:"merchant_id" json:"-" faker:"objectId"`
	ProjectId       primitive.ObjectID              `bson:"project_id" json:"project_id" faker:"objectId"`
	Pricing         string                          `bson:"pricing" json:"pricing"`
	BillingType     string                          `bson:"billing_type" json:"billing_type"`
}

func (m *productMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.Product)

	out := &MgoProduct{
		Object:          in.Object,
		Type:            in.Type,
		Sku:             in.Sku,
		DefaultCurrency: in.DefaultCurrency,
		Enabled:         in.Enabled,
		Description:     in.Description,
		LongDescription: in.LongDescription,
		Images:          in.Images,
		Url:             in.Url,
		Metadata:        in.Metadata,
		Deleted:         in.Deleted,
		Pricing:         in.Pricing,
		BillingType:     in.BillingType,
	}

	merchantOid, err := primitive.ObjectIDFromHex(in.MerchantId)

	if err != nil {
		return nil, err
	}

	out.MerchantId = merchantOid

	projectOid, err := primitive.ObjectIDFromHex(in.ProjectId)

	if err != nil {
		return nil, err
	}

	out.ProjectId = projectOid

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

	out.Name = []*billingpb.I18NTextSearchable{}
	for k, v := range in.Name {
		out.Name = append(out.Name, &billingpb.I18NTextSearchable{Lang: k, Value: v})
	}

	for _, price := range in.Prices {
		out.Prices = append(out.Prices, &billingpb.ProductPrice{
			Currency:          price.Currency,
			Region:            price.Region,
			Amount:            tools.FormatAmount(price.Amount),
			IsVirtualCurrency: price.IsVirtualCurrency,
		})
	}

	return out, nil
}

func (m *productMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoProduct)

	out := &billingpb.Product{
		Id:              in.Id.Hex(),
		Object:          in.Object,
		Type:            in.Type,
		Sku:             in.Sku,
		DefaultCurrency: in.DefaultCurrency,
		Enabled:         in.Enabled,
		Prices:          in.Prices,
		Description:     in.Description,
		LongDescription: in.LongDescription,
		Images:          in.Images,
		Url:             in.Url,
		Metadata:        in.Metadata,
		Deleted:         in.Deleted,
		MerchantId:      in.MerchantId.Hex(),
		ProjectId:       in.ProjectId.Hex(),
		Pricing:         in.Pricing,
		BillingType:     in.BillingType,
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt)

	if err != nil {
		return nil, err
	}

	out.Name = map[string]string{}
	for _, i := range in.Name {
		out.Name[i.Lang] = i.Value
	}

	return out, nil
}
