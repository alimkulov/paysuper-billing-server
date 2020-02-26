package repository

import (
	"context"
	pkg2 "github.com/paysuper/paysuper-billing-server/internal/pkg"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"time"
)

const (
	collectionAccountingEntry = "accounting_entry"
)

type accountingEntryRepository repository

// NewAccountingEntryRepository create and return an object for working with the accounting entry repository.
// The returned object implements the AccountingEntryRepositoryInterface interface.
func NewAccountingEntryRepository(db mongodb.SourceInterface) AccountingEntryRepositoryInterface {
	s := &accountingEntryRepository{db: db, mapper: models.NewAccountingEntryMapper()}
	return s
}

func (r *accountingEntryRepository) MultipleInsert(ctx context.Context, objs []*billingpb.AccountingEntry) error {
	c := make([]interface{}, len(objs))
	for i, v := range objs {
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

	_, err := r.db.Collection(collectionAccountingEntry).InsertMany(ctx, c)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, c),
		)
		return err
	}

	return nil
}

func (r *accountingEntryRepository) GetById(ctx context.Context, id string) (*billingpb.AccountingEntry, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	var mgo = models.MgoAccountingEntry{}
	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionAccountingEntry).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
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

	return obj.(*billingpb.AccountingEntry), nil
}

func (r *accountingEntryRepository) GetByObjectSource(ctx context.Context, object, sourceId, sourceType string) (*billingpb.AccountingEntry, error) {
	objId, err := primitive.ObjectIDFromHex(sourceId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, sourceId),
		)
		return nil, err
	}

	query := bson.M{
		"object":      object,
		"source.id":   objId,
		"source.type": sourceType,
	}

	var mgo = models.MgoAccountingEntry{}
	err = r.db.Collection(collectionAccountingEntry).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
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

	return obj.(*billingpb.AccountingEntry), nil
}

func (r *accountingEntryRepository) ApplyObjectSource(
	ctx context.Context, object, tp, sourceId, sourceType string, source *billingpb.AccountingEntry,
) (*billingpb.AccountingEntry, error) {
	objId, err := primitive.ObjectIDFromHex(sourceId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, sourceId),
		)
		return nil, err
	}

	query := bson.M{
		"object":      object,
		"type":        tp,
		"source.id":   objId,
		"source.type": sourceType,
	}

	obj, err := r.mapper.MapObjectToMgo(source)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, source),
		)
		return nil, err
	}

	var mgo = obj.(*models.MgoAccountingEntry)
	err = r.db.Collection(collectionAccountingEntry).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	ae, err := r.mapper.MapMgoToObject(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return ae.(*billingpb.AccountingEntry), nil
}

func (r *accountingEntryRepository) FindBySource(ctx context.Context, sourceId, sourceType string) ([]*billingpb.AccountingEntry, error) {
	objId, err := primitive.ObjectIDFromHex(sourceId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, sourceId),
		)
		return nil, err
	}

	query := bson.M{
		"source.id":   objId,
		"source.type": sourceType,
	}

	cursor, err := r.db.Collection(collectionAccountingEntry).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoAccountingEntry
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.AccountingEntry, len(list))

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
		objs[i] = v.(*billingpb.AccountingEntry)
	}

	return objs, nil
}

func (r *accountingEntryRepository) GetByTypeWithTaxes(ctx context.Context, types []string) ([]*billingpb.AccountingEntry, error) {
	query := bson.M{
		"type":   bson.M{"$in": types},
		"amount": bson.M{"$gt": 0},
	}

	opts := options.Find().
		SetSort(mongodb.ToSortOption([]string{"source.type", "source.id", "-type"}))

	cursor, err := r.db.Collection(collectionAccountingEntry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoAccountingEntry
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.AccountingEntry, len(list))

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
		objs[i] = v.(*billingpb.AccountingEntry)
	}

	return objs, nil
}

func (r *accountingEntryRepository) FindByTypeCountryDates(
	ctx context.Context, country string, types []string, dateFrom, dateTo time.Time,
) ([]*billingpb.AccountingEntry, error) {
	query := bson.M{
		"created_at": bson.M{
			"$gte": dateFrom,
			"$lte": dateTo,
		},
		"country": country,
		"type":    bson.M{"$in": types},
	}

	cursor, err := r.db.Collection(collectionAccountingEntry).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoAccountingEntry
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.AccountingEntry, len(list))

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
		objs[i] = v.(*billingpb.AccountingEntry)
	}

	return objs, nil
}

func (r *accountingEntryRepository) GetDistinctBySourceId(ctx context.Context) ([]string, error) {
	res, err := r.db.Collection(collectionAccountingEntry).Distinct(ctx, "source.id", bson.M{})

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
		)
		return nil, err
	}

	var ids []string

	for _, id := range res {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

func (r *accountingEntryRepository) GetCorrectionsForRoyaltyReport(
	ctx context.Context, merchantId, currency string, from, to time.Time,
) ([]*billingpb.AccountingEntry, error) {
	id, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	query := bson.M{
		"merchant_id": id,
		"currency":    currency,
		"created_at":  bson.M{"$gte": from, "$lte": to},
		"type":        pkg.AccountingEntryTypeMerchantRoyaltyCorrection,
	}

	sorts := bson.M{"created_at": 1}
	opts := options.Find()
	opts.SetSort(sorts)
	cursor, err := r.db.Collection(collectionAccountingEntry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return nil, err
	}

	var list []*models.MgoAccountingEntry
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return nil, err
	}

	objs := make([]*billingpb.AccountingEntry, len(list))

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
		objs[i] = v.(*billingpb.AccountingEntry)
	}

	return objs, nil
}

func (r *accountingEntryRepository) GetRollingReservesForRoyaltyReport(
	ctx context.Context, merchantId, currency string, types []string, from, to time.Time,
) ([]*billingpb.AccountingEntry, error) {
	id, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	query := bson.M{
		"merchant_id": id,
		"currency":    currency,
		"created_at":  bson.M{"$gte": from, "$lte": to},
		"type":        bson.M{"$in": types},
	}

	sorts := bson.M{"created_at": 1}
	opts := options.Find()
	opts.SetSort(sorts)
	cursor, err := r.db.Collection(collectionAccountingEntry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return nil, err
	}

	var list []*models.MgoAccountingEntry
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return nil, err
	}

	objs := make([]*billingpb.AccountingEntry, len(list))

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
		objs[i] = v.(*billingpb.AccountingEntry)
	}

	return objs, nil
}

func (r *accountingEntryRepository) GetRollingReserveForBalance(
	ctx context.Context, merchantId, currency string, types []string, createdAt time.Time,
) ([]*pkg2.ReserveQueryResItem, error) {
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	matchQuery := bson.M{
		"merchant_id": merchantOid,
		"currency":    currency,
		"type":        bson.M{"$in": types},
	}

	if !createdAt.IsZero() {
		matchQuery["created_at"] = bson.M{"$gt": createdAt}
	}

	query := []bson.M{
		{
			"$match": matchQuery,
		},
		{
			"$group": bson.M{"_id": "$type", "amount": bson.M{"$sum": "$amount"}},
		},
	}

	var items []*pkg2.ReserveQueryResItem
	cursor, err := r.db.Collection(collectionAccountingEntry).Aggregate(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	err = cursor.All(ctx, &items)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	return items, nil
}

func (r *accountingEntryRepository) BulkWrite(ctx context.Context, list []*billingpb.AccountingEntry) error {
	var operations []mongo.WriteModel

	for _, ae := range list {
		oid, err := primitive.ObjectIDFromHex(ae.Id)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
				zap.String(pkg.ErrorDatabaseFieldQuery, ae.Id),
			)
			return err
		}

		mgo, err := r.mapper.MapObjectToMgo(ae)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, ae),
			)
			return err
		}

		operation := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": oid}).
			SetUpdate(mgo)
		operations = append(operations, operation)
	}

	_, err := r.db.Collection(collectionAccountingEntry).BulkWrite(ctx, operations)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAccountingEntry),
		)
		return err
	}

	return nil
}
