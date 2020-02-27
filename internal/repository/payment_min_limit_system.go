package repository

import (
	"context"
	"fmt"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	tools "github.com/paysuper/paysuper-tools/number"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionPaymentMinLimitSystem = "payment_min_limit_system"

	cacheKeyPaymentMinLimitSystem    = "payment_min_limit_system:currency:%s"
	cacheKeyAllPaymentMinLimitSystem = "payment_min_limit_system:all"
)

type paymentMinLimitSystemRepository repository

// NewPaymentMinLimitSystemRepository create and return an object for working with the payment min limit repository.
// The returned object implements the PaymentMinLimitSystemRepositoryInterface interface.
func NewPaymentMinLimitSystemRepository(db mongodb.SourceInterface, cache database.CacheInterface) PaymentMinLimitSystemRepositoryInterface {
	s := &paymentMinLimitSystemRepository{db: db, cache: cache}
	return s
}

func (r *paymentMinLimitSystemRepository) MultipleInsert(
	ctx context.Context,
	pmlsArray []*billingpb.PaymentMinLimitSystem,
) error {
	c := make([]interface{}, len(pmlsArray))
	for i, v := range pmlsArray {
		v.Amount = tools.FormatAmount(v.Amount)
		mgo, err := r.mapper.MapObjectToMgo(v)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, v),
			)
		}

		c[i] = mgo
	}

	_, err := r.db.Collection(collectionPaymentMinLimitSystem).InsertMany(ctx, c)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentMinLimitSystem),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, c),
		)
		return err
	}

	return nil
}

func (r *paymentMinLimitSystemRepository) Upsert(ctx context.Context, pmls *billingpb.PaymentMinLimitSystem) error {
	pmls.Amount = tools.FormatAmount(pmls.Amount)
	oid, err := primitive.ObjectIDFromHex(pmls.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentChannelCostSystem),
			zap.String(pkg.ErrorDatabaseFieldQuery, pmls.Id),
		)
	}

	mgo, err := r.mapper.MapObjectToMgo(pmls)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, pmls),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	opts := options.Replace().SetUpsert(true)
	_, err = r.db.Collection(collectionPaymentMinLimitSystem).ReplaceOne(ctx, filter, mgo, opts)

	if err != nil {
		zap.S().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentMinLimitSystem),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldDocument, pmls),
		)
		return err
	}

	key := fmt.Sprintf(cacheKeyPaymentMinLimitSystem, pmls.Currency)
	err = r.cache.Set(key, pmls, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, pmls),
		)
		return err
	}

	err = r.cache.Delete(cacheKeyAllPaymentMinLimitSystem)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "DELETE"),
			zap.String(pkg.ErrorCacheFieldKey, cacheKeyAllPaymentMinLimitSystem),
		)
	}

	return nil
}

func (r *paymentMinLimitSystemRepository) GetByCurrency(
	ctx context.Context,
	currency string,
) (*billingpb.PaymentMinLimitSystem, error) {
	pmls := &billingpb.PaymentMinLimitSystem{}
	key := fmt.Sprintf(cacheKeyPaymentMinLimitSystem, currency)

	if err := r.cache.Get(key, &pmls); err == nil {
		return pmls, nil
	}

	var mgo = models.MgoPriceGroup{}
	query := bson.M{"currency": currency}
	err := r.db.Collection(collectionPaymentMinLimitSystem).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentMinLimitSystem),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := r.mapper.MapMgoToObject(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	pmls = obj.(*billingpb.PaymentMinLimitSystem)

	err = r.cache.Set(key, pmls, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, pmls),
		)
	}

	return pmls, nil
}

func (r *paymentMinLimitSystemRepository) GetAll(ctx context.Context) ([]*billingpb.PaymentMinLimitSystem, error) {
	result := []*billingpb.PaymentMinLimitSystem{}

	if err := r.cache.Get(cacheKeyAllPaymentMinLimitSystem, &result); err == nil {
		return result, nil
	}

	cursor, err := r.db.Collection(collectionPaymentMinLimitSystem).Find(ctx, bson.M{})

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentMinLimitSystem),
			zap.Any(pkg.ErrorDatabaseFieldQuery, nil),
		)
		return nil, err
	}

	var list []*models.MgoPaymentMinLimitSystem
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaymentMinLimitSystem),
			zap.Any(pkg.ErrorDatabaseFieldQuery, nil),
		)
		return nil, err
	}

	objs := make([]*billingpb.PaymentMinLimitSystem, len(list))

	for i, obj := range list {
		v, err := r.mapper.MapMgoToObject(obj)
		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
			)
			return nil, err
		}
		objs[i] = v.(*billingpb.PaymentMinLimitSystem)
	}

	err = r.cache.Set(cacheKeyAllPaymentMinLimitSystem, objs, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, collectionPaymentMinLimitSystem),
			zap.Any(pkg.ErrorCacheFieldData, objs),
		)
	}

	return objs, nil
}
