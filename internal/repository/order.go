package repository

import (
	"context"
	"github.com/paysuper/paysuper-billing-server/internal/helper"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"time"
)

const (
	// CollectionOrder is name of table for collection the order.
	CollectionOrder = "order"
)

type orderRepository repository

// NewOrderRepository create and return an object for working with the order repository.
// The returned object implements the OrderRepositoryInterface interface.
func NewOrderRepository(db mongodb.SourceInterface) OrderRepositoryInterface {
	s := &orderRepository{db: db, mapper: models.NewOrderMapper()}
	return s
}

func (h *orderRepository) Insert(ctx context.Context, order *billingpb.Order) error {
	mgo, err := h.mapper.MapObjectToMgo(order)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, order),
		)
		return err
	}

	_, err = h.db.Collection(CollectionOrder).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, order),
		)
		return err
	}

	return nil
}

func (h *orderRepository) Update(ctx context.Context, order *billingpb.Order) error {
	oid, err := primitive.ObjectIDFromHex(order.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.String(pkg.ErrorDatabaseFieldQuery, order.Id),
		)
		return err
	}

	mgo, err := h.mapper.MapObjectToMgo(order)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, order),
		)
		return err
	}

	_, err = h.db.Collection(CollectionOrder).ReplaceOne(ctx, bson.M{"_id": oid}, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, order),
		)
		return err
	}

	return nil
}

func (h *orderRepository) GetByUuid(ctx context.Context, uuid string) (*billingpb.Order, error) {
	query := bson.M{"uuid": uuid}

	mgo := &models.MgoOrder{}
	err := h.db.Collection(CollectionOrder).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := h.mapper.MapMgoToObject(mgo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return obj.(*billingpb.Order), nil
}

func (h *orderRepository) GetById(ctx context.Context, id string) (*billingpb.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	query := bson.M{"_id": oid}
	mgo := &models.MgoOrder{}
	err = h.db.Collection(CollectionOrder).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := h.mapper.MapMgoToObject(mgo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return obj.(*billingpb.Order), nil
}

func (h *orderRepository) GetByRefundReceiptNumber(ctx context.Context, id string) (*billingpb.Order, error) {
	query := bson.M{"refund.receipt_number": id}
	mgo := &models.MgoOrder{}

	err := h.db.Collection(CollectionOrder).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := h.mapper.MapMgoToObject(mgo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return obj.(*billingpb.Order), nil
}

func (h *orderRepository) GetByProjectOrderId(ctx context.Context, projectId, projectOrderId string) (*billingpb.Order, error) {
	id, err := primitive.ObjectIDFromHex(projectId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
		)
		return nil, err
	}

	mgo := &models.MgoOrder{}

	query := bson.M{"project._id": id, "project_order_id": projectOrderId}
	err = h.db.Collection(CollectionOrder).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := h.mapper.MapMgoToObject(mgo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return obj.(*billingpb.Order), nil
}

func (h *orderRepository) UpdateOrderView(ctx context.Context, ids []string) error {
	defer helper.TimeTrack(time.Now(), "updateOrderView")

	idsHex := []primitive.ObjectID{}

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			continue
		}

		idsHex = append(idsHex, oid)
	}

	var match bson.M

	if len(idsHex) == 1 {
		match = bson.M{
			"$match": bson.M{
				"_id": idsHex[0],
			},
		}
	} else {
		match = bson.M{
			"$match": bson.M{
				"_id": bson.M{"$in": idsHex},
			},
		}
	}

	orderViewQuery := []bson.M{
		match,
		{
			"$lookup": bson.M{
				"from": "merchant",
				"let": bson.M{
					"merchant_id": "$project.merchant_id",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []string{
									"$_id",
									"$$merchant_id",
								},
							},
						},
					},
					{
						"$project": bson.M{
							"company_name":     "$company.name",
							"agreement_number": "$agreement_number",
							"_id":              0,
						},
					},
				},
				"as": "merchant_info",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$merchant_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_gross_revenue",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$local_amount",
							"currency": "$local_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_gross_revenue",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$original_amount",
							"currency": "$original_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue_origin",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue_origin",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_gross_revenue",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_tax_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_tax_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$local_amount",
							"currency": "$local_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_tax_fee_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_tax_fee_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$original_amount",
							"currency": "$original_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_tax_fee_origin",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_tax_fee_origin",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "central_bank_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_tax_fee_current_exchange_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_tax_fee_current_exchange_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "ps_gross_revenue_fx",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue_fx",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue_fx",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "ps_gross_revenue_fx_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue_fx_tax_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue_fx_tax_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_tax_fee_cost_value",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "tax_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$tax_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_tax_fee_central_bank_fx",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "tax_fee_currency_exchange_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$tax_fee_currency_exchange_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "ps_method_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "method_fee_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$method_fee_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_method_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "method_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$method_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_method_fee_cost_value",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_fee_tariff_self_cost",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_fee_tariff_self_cost",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_method_fixed_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "method_fixed_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$method_fixed_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_merchant_method_fixed_fee_cost_value",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_fixed_fee_tariff_self_cost",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_fixed_fee_tariff_self_cost",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_ps_fixed_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_fixed_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_fixed_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$local_amount",
							"currency": "$local_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_gross_revenue_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_gross_revenue_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$original_amount",
							"currency": "$original_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_gross_revenue_origin",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_gross_revenue_origin",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_gross_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_gross_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$local_amount",
							"currency": "$local_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_tax_fee_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_tax_fee_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   "$original_amount",
							"currency": "$original_currency",
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_tax_fee_origin",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_tax_fee_origin",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_tax_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_tax_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_refund_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_refund_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "real_refund_fixed_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "method_refund_fixed_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$method_refund_fixed_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_refund",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_gross_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_gross_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_refund_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "method_refund_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$method_refund_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_refund_fixed_fee_cost_value",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_refund_fixed_fee_tariff_self_cost",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_refund_fixed_fee_tariff_self_cost",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "merchant_refund_fixed_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "merchant_refund_fixed_fee_tariff",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$merchant_refund_fixed_fee_tariff",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "reverse_tax_fee",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_tax_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_tax_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "reverse_tax_fee_delta",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_tax_fee_currency_exchange_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_tax_fee_currency_exchange_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": "ps_reverse_tax_fee_delta",
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_refund_tax_fee_currency_exchange_fee",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_refund_tax_fee_currency_exchange_fee",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"real_tax_fee",
									"central_bank_tax_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": "$amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_tax_fee_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_tax_fee_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_tax_fee_cost_value",
									"merchant_tax_fee_central_bank_fx",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": "$amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "tax_fee_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$tax_fee_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"reverse_tax_fee",
									"reverse_tax_fee_delta",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": "$amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_tax_fee_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_tax_fee_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"ps_method_fee",
									"merchant_ps_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": "$amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "fees_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$fees_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"ps_method_fee",
									"merchant_ps_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$local_currency",
							"amount": bson.M{
								"$sum": "$local_amount",
							},
							"currency": bson.M{
								"$first": "$local_currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "fees_total_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$fees_total_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund_fee",
									"merchant_refund_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": "$amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_fees_total",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_fees_total",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund_fee",
									"merchant_refund_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$local_currency",
							"amount": bson.M{
								"$sum": "$local_amount",
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_fees_total_local",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_fees_total_local",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"ps_gross_revenue_fx",
									"ps_gross_revenue_fx_tax_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"ps_gross_revenue_fx_tax_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "payment_gross_revenue_fx_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$payment_gross_revenue_fx_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"real_gross_revenue",
									"ps_gross_revenue_fx",
								},
							},
						},
					},

					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"ps_gross_revenue_fx",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "gross_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$gross_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_method_fee",
									"merchant_method_fee_cost_value",
								},
							},
						},
					},

					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"merchant_method_fee_cost_value",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_fee_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_fee_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_method_fixed_fee",
									"real_merchant_method_fixed_fee",
								},
							},
						},
					},

					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_merchant_method_fixed_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_fixed_fee_tariff_fx_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_fixed_fee_tariff_fx_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"real_merchant_method_fixed_fee",
									"real_merchant_method_fixed_fee_cost_value",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_merchant_method_fixed_fee_cost_value",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_fixed_fee_tariff_total_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_fixed_fee_tariff_total_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_ps_fixed_fee",
									"real_merchant_ps_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_merchant_ps_fixed_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_fixed_fee_fx_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_fixed_fee_fx_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund",
									"real_refund",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_refund",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_gross_revenue_fx",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_gross_revenue_fx",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},

							"type": bson.M{
								"$in": []string{
									"merchant_refund_fee",
									"real_refund_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_refund_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_refund_fee_tariff_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_refund_fee_tariff_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund_fixed_fee",
									"real_refund_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"real_refund_fixed_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_refund_fixed_fee_tariff_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_refund_fixed_fee_tariff_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"real_gross_revenue",
									"merchant_tax_fee_central_bank_fx",
									"ps_gross_revenue_fx",
									"merchant_tax_fee_cost_value",
									"ps_method_fee",
									"merchant_ps_fixed_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$in": []interface{}{"$type", []string{
												"ps_gross_revenue_fx",
												"merchant_tax_fee_central_bank_fx",
												"merchant_tax_fee_cost_value",
												"ps_method_fee",
												"merchant_ps_fixed_fee",
											},
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "net_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$net_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"ps_method_fee",
									"merchant_ps_fixed_fee",
									"merchant_method_fee_cost_value",
									"real_merchant_method_fixed_fee_cost_value",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$in": []interface{}{"$type", []string{
												"merchant_method_fee_cost_value",
												"real_merchant_method_fixed_fee_cost_value",
											},
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_method_total_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_method_total_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"ps_gross_revenue_fx",
									"ps_method_fee",
									"merchant_ps_fixed_fee",
									"central_bank_tax_fee",
									"ps_gross_revenue_fx_tax_fee",
									"merchant_method_fee_cost_value",
									"real_merchant_method_fixed_fee_cost_value",
								},
							},
						},
					},

					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$in": []interface{}{"$type", []string{
												"central_bank_tax_fee",
												"ps_gross_revenue_fx_tax_fee",
												"merchant_method_fee_cost_value",
												"real_merchant_method_fixed_fee_cost_value",
											},
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_total_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_total_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund",
									"merchant_refund_fee",
									"merchant_refund_fixed_fee",
									"reverse_tax_fee_delta",
									"reverse_tax_fee",
								},
							},
						},
					},
					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$eq": []string{
												"$type",
												"reverse_tax_fee",
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "refund_reverse_revenue",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$refund_reverse_revenue",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from": "accounting_entry",
				"let": bson.M{
					"order_id":    "$_id",
					"object_type": "$type",
				},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$and": []interface{}{
									bson.M{
										"$eq": []string{
											"$source.id",
											"$$order_id",
										}},
									bson.M{
										"$eq": []string{
											"$source.type",
											"$$object_type",
										}},
								},
							},
							"type": bson.M{
								"$in": []string{
									"merchant_refund_fee",
									"merchant_refund_fixed_fee",
									"ps_reverse_tax_fee_delta",
									"real_refund_fixed_fee",
									"real_refund_fee",
								},
							},
						},
					},

					{
						"$group": bson.M{
							"_id": "$currency",
							"amount": bson.M{
								"$sum": bson.M{
									"$cond": []interface{}{
										bson.M{
											"$in": []interface{}{"$type", []string{
												"real_refund_fixed_fee",
												"real_refund_fee",
											},
											},
										},
										bson.M{
											"$subtract": []interface{}{
												0,
												"$amount",
											},
										},
										"$amount",
									},
								},
							},
							"currency": bson.M{
								"$first": "$currency",
							},
						},
					},
					{
						"$project": bson.M{
							"amount":   1,
							"currency": 1,
							"_id":      0,
						},
					},
				},
				"as": "paysuper_refund_total_profit",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$paysuper_refund_total_profit",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$addFields": bson.M{
				"order_charge": bson.M{
					"amount":   "$charge_amount",
					"currency": "$charge_currency",
				},
			},
		},
		{
			"$addFields": bson.M{
				"order_charge_before_vat": bson.M{
					"amount": bson.M{
						"$cond": []interface{}{
							bson.M{
								"$eq": []string{
									"$type",
									"order",
								},
							},
							bson.M{"$subtract": []interface{}{"$charge_amount", "$payment_tax_fee_origin.amount"}},
							bson.M{"$subtract": []interface{}{"$charge_amount", "$payment_refund_tax_fee_origin.amount"}},
						},
					},
					"currency": "$charge_currency",
				},
			},
		},
		{
			"$project": bson.M{
				"_id":                  1,
				"uuid":                 1,
				"pm_order_id":          1,
				"project_order_id":     1,
				"project":              1,
				"created_at":           1,
				"pm_order_close_date":  1,
				"total_payment_amount": 1,
				"amount_before_vat":    "$private_amount",
				"currency":             1,
				"user":                 1,
				"billing_address":      1,
				"payment_method":       1,
				"country_code":         1,
				"merchant_id":          "$project.merchant_id",
				"status":               1,
				"tax_rate":             "$tax.rate",
				"merchant_info":        1,
				"locale": bson.M{
					"$cond": []interface{}{
						bson.M{
							"$ne": []interface{}{"$user", nil},
						},
						"$user.locale",
						"",
					},
				},
				"type":                                              1,
				"royalty_report_id":                                 1,
				"is_vat_deduction":                                  1,
				"payment_gross_revenue_local":                       1,
				"payment_gross_revenue_origin":                      1,
				"payment_gross_revenue":                             1,
				"payment_tax_fee":                                   1,
				"payment_tax_fee_local":                             1,
				"payment_tax_fee_origin":                            1,
				"payment_tax_fee_current_exchange_fee":              1,
				"payment_tax_fee_total":                             1,
				"payment_gross_revenue_fx":                          1,
				"payment_gross_revenue_fx_tax_fee":                  1,
				"payment_gross_revenue_fx_profit":                   1,
				"gross_revenue":                                     1,
				"tax_fee":                                           1,
				"tax_fee_currency_exchange_fee":                     1,
				"tax_fee_total":                                     1,
				"method_fee_total":                                  1,
				"method_fee_tariff":                                 1,
				"paysuper_method_fee_tariff_self_cost":              1,
				"paysuper_method_fee_profit":                        1,
				"method_fixed_fee_tariff":                           1,
				"paysuper_method_fixed_fee_tariff_fx_profit":        1,
				"paysuper_method_fixed_fee_tariff_self_cost":        1,
				"paysuper_method_fixed_fee_tariff_total_profit":     1,
				"paysuper_fixed_fee":                                1,
				"paysuper_fixed_fee_fx_profit":                      1,
				"fees_total":                                        1,
				"fees_total_local":                                  1,
				"net_revenue":                                       1,
				"paysuper_method_total_profit":                      1,
				"paysuper_total_profit":                             1,
				"payment_refund_gross_revenue_local":                1,
				"payment_refund_gross_revenue_origin":               1,
				"payment_refund_gross_revenue":                      1,
				"payment_refund_tax_fee":                            1,
				"payment_refund_tax_fee_local":                      1,
				"payment_refund_tax_fee_origin":                     1,
				"payment_refund_fee_tariff":                         1,
				"method_refund_fixed_fee_tariff":                    1,
				"refund_gross_revenue":                              1,
				"refund_gross_revenue_fx":                           1,
				"method_refund_fee_tariff":                          1,
				"paysuper_method_refund_fee_tariff_profit":          1,
				"paysuper_method_refund_fixed_fee_tariff_self_cost": 1,
				"merchant_refund_fixed_fee_tariff":                  1,
				"paysuper_method_refund_fixed_fee_tariff_profit":    1,
				"refund_tax_fee":                                    1,
				"refund_tax_fee_currency_exchange_fee":              1,
				"paysuper_refund_tax_fee_currency_exchange_fee":     1,
				"refund_tax_fee_total":                              1,
				"refund_reverse_revenue":                            1,
				"refund_fees_total":                                 1,
				"refund_fees_total_local":                           1,
				"paysuper_refund_total_profit":                      1,
				"issuer":                                            1,
				"items":                                             1,
				"parent_order":                                      1,
				"refund":                                            1,
				"cancellation":                                      1,
				"mcc_code":                                          1,
				"operating_company_id":                              1,
				"is_high_risk":                                      1,
				"payment_ip_country":                                1,
				"is_ip_country_mismatch_bin":                        1,
				"order_charge":                                      1,
				"order_charge_before_vat":                           1,
				"billing_country_changed_by_user":                   1,
				"refund_allowed":                                    "$is_refund_allowed",
				"vat_payer":                                         1,
				"is_production":                                     1,
				"merchant_payout_currency": bson.M{
					"$ifNull": []interface{}{"$net_revenue.currency", "$refund_reverse_revenue.currency"},
				},
			},
		},
		{
			"$project": bson.M{
				"payment_method.params":            0,
				"payment_method.payment_system_id": 0,
			},
		},
		{
			"$merge": bson.M{
				"into":        "order_view",
				"whenMatched": "replace",
			},
		},
	}

	cursor, err := h.db.Collection(CollectionOrder).Aggregate(ctx, orderViewQuery)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
		)
		return err
	}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			zap.L().Error(
				pkg.ErrorQueryCursorCloseFailed,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
			)
		}
	}()

	cursor.Next(ctx)
	err = cursor.Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionOrder),
		)
		return err
	}

	return nil
}
