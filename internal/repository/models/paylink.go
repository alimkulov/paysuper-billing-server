package models

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type payLinkMapper struct{}

func NewPayLinkMapperMapper() Mapper {
	return &payLinkMapper{}
}

type MgoPaylink struct {
	Id                   primitive.ObjectID `bson:"_id"`
	Object               string             `bson:"object"`
	Products             []string           `bson:"products"`
	ExpiresAt            time.Time          `bson:"expires_at"`
	CreatedAt            time.Time          `bson:"created_at"`
	UpdatedAt            time.Time          `bson:"updated_at"`
	MerchantId           primitive.ObjectID `bson:"merchant_id"`
	ProjectId            primitive.ObjectID `bson:"project_id"`
	Name                 string             `bson:"name"`
	ProductsType         string             `bson:"products_type"`
	IsExpired            bool               `bson:"is_expired"`
	Visits               int32              `bson:"visits"`
	NoExpiryDate         bool               `bson:"no_expiry_date"`
	TotalTransactions    int32              `bson:"total_transactions"`
	SalesCount           int32              `bson:"sales_count"`
	ReturnsCount         int32              `bson:"returns_count"`
	Conversion           float64            `bson:"conversion"`
	GrossSalesAmount     float64            `bson:"gross_sales_amount"`
	GrossReturnsAmount   float64            `bson:"gross_returns_amount"`
	GrossTotalAmount     float64            `bson:"gross_total_amount"`
	TransactionsCurrency string             `bson:"transactions_currency"`
	Deleted              bool               `bson:"deleted"`
}

func (m *payLinkMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	in := obj.(*billingpb.Paylink)

	out := &MgoPaylink{
		Object:               in.Object,
		Products:             in.Products,
		Name:                 in.Name,
		ProductsType:         in.ProductsType,
		Deleted:              in.Deleted,
		NoExpiryDate:         in.NoExpiryDate,
		TotalTransactions:    in.TotalTransactions,
		SalesCount:           in.SalesCount,
		ReturnsCount:         in.ReturnsCount,
		Conversion:           in.Conversion,
		GrossSalesAmount:     in.GrossSalesAmount,
		GrossReturnsAmount:   in.GrossReturnsAmount,
		GrossTotalAmount:     in.GrossTotalAmount,
		TransactionsCurrency: in.TransactionsCurrency,
		Visits:               in.Visits,
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

	if in.ExpiresAt != nil {
		t, err := ptypes.Timestamp(in.ExpiresAt)

		if err != nil {
			return nil, err
		}

		out.ExpiresAt = t
	}

	out.IsExpired = in.IsPaylinkExpired()

	return out, nil
}

func (m *payLinkMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	var err error
	in := obj.(*MgoPaylink)

	out := &billingpb.Paylink{
		Id:                   in.Id.Hex(),
		MerchantId:           in.MerchantId.Hex(),
		ProjectId:            in.ProjectId.Hex(),
		Object:               in.Object,
		Products:             in.Products,
		Name:                 in.Name,
		ProductsType:         in.ProductsType,
		Deleted:              in.Deleted,
		NoExpiryDate:         in.NoExpiryDate,
		TotalTransactions:    in.TotalTransactions,
		SalesCount:           in.SalesCount,
		ReturnsCount:         in.ReturnsCount,
		Conversion:           in.Conversion,
		GrossSalesAmount:     in.GrossSalesAmount,
		GrossReturnsAmount:   in.GrossReturnsAmount,
		GrossTotalAmount:     in.GrossTotalAmount,
		TransactionsCurrency: in.TransactionsCurrency,
		Visits:               in.Visits,
	}

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)
	if err != nil {
		return nil, err
	}

	out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt)
	if err != nil {
		return nil, err
	}

	out.ExpiresAt, err = ptypes.TimestampProto(in.ExpiresAt)
	if err != nil {
		return nil, err
	}

	out.IsExpired = out.IsPaylinkExpired()

	return out, nil
}
