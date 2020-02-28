package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	pkg2 "github.com/paysuper/paysuper-billing-server/internal/pkg"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"time"
)

const (
	CollectionRoyaltyReport        = "royalty_report"
	CollectionRoyaltyReportChanges = "royalty_report_changes"

	cacheKeyRoyaltyReport = "royalty_report:id:%s"
)

type royaltyReportRepository repository

// NewRoyaltyReportRepository create and return an object for working with the royalty report repository.
// The returned object implements the RoyaltyReportRepositoryInterface interface.
func NewRoyaltyReportRepository(db mongodb.SourceInterface, cache database.CacheInterface) RoyaltyReportRepositoryInterface {
	s := &royaltyReportRepository{db: db, cache: cache, mapper: models.NewRoyaltyReportMapper()}
	return s
}

func (r *royaltyReportRepository) GetNonPayoutReports(
	ctx context.Context,
	merchantId, currency string,
) (result []*billingpb.RoyaltyReport, err error) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	royaltyReportsStatusActive := []string{
		billingpb.RoyaltyReportStatusPending,
		billingpb.RoyaltyReportStatusAccepted,
		billingpb.RoyaltyReportStatusDispute,
	}
	query := bson.M{
		"merchant_id":        oid,
		"currency":           currency,
		"status":             bson.M{"$in": royaltyReportsStatusActive},
		"payout_document_id": "",
	}

	sorts := bson.M{"period_from": 1}
	opts := options.Find().SetSort(sorts)
	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	result = objs

	return
}

func (r *royaltyReportRepository) GetByPayoutId(ctx context.Context, payoutId string) (result []*billingpb.RoyaltyReport, err error) {
	query := bson.M{
		"payout_document_id": payoutId,
	}

	sorts := bson.M{"period_from": 1}
	opts := options.Find().SetSort(sorts)
	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSorts, sorts),
		)
		return
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	result = objs

	return
}

func (r *royaltyReportRepository) GetAll(ctx context.Context) ([]*billingpb.RoyaltyReport, error) {
	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, bson.M{})

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
		)
		return nil, err
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
		)
		return nil, err
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	return objs, nil
}

func (r *royaltyReportRepository) GetByPeriod(ctx context.Context, from, to time.Time) ([]*billingpb.RoyaltyReport, error) {
	query := bson.M{"period_from": from, "period_to": to}
	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	return objs, nil
}

func (r *royaltyReportRepository) GetByAcceptedExpireWithStatus(
	ctx context.Context, date time.Time, status string,
) ([]*billingpb.RoyaltyReport, error) {
	query := bson.M{
		"accept_expire_at": bson.M{"$lte": date},
		"status":           status,
	}

	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	return objs, nil
}

func (r *royaltyReportRepository) GetRoyaltyHistoryById(ctx context.Context, id string) ([]*billingpb.RoyaltyReportChanges, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	query := bson.M{"royalty_report_id": oid}
	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := r.db.Collection(CollectionRoyaltyReportChanges).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoRoyaltyReportChanges
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.RoyaltyReportChanges, len(list))

	for i, obj := range list {
		v, err := models.NewRoyaltyReportChangesMapper().MapMgoToObject(obj)
		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
			)
			return nil, err
		}
		objs[i] = v.(*billingpb.RoyaltyReportChanges)
	}

	return objs, nil
}

func (r *royaltyReportRepository) FindByMerchantStatusDates(
	ctx context.Context, merchantId string, status []string, dateFrom, dateTo, offset, limit int64,
) ([]*billingpb.RoyaltyReport, error) {
	var err error
	query := bson.M{}

	if merchantId != "" {
		query["merchant_id"], err = primitive.ObjectIDFromHex(merchantId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return nil, err
		}
	}

	if len(status) > 0 {
		query["status"] = bson.M{"$in": status}
	}

	if dateFrom > 0 || dateTo > 0 {
		date := bson.M{}
		if dateFrom > 0 {
			date["$gte"] = time.Unix(dateFrom, 0)
		}
		if dateTo > 0 {
			date["$lte"] = time.Unix(dateTo, 0)
		}
		query["created_at"] = date
	}

	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset)
	cursor, err := r.db.Collection(CollectionRoyaltyReport).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoRoyaltyReport
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.RoyaltyReport, len(list))

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
		objs[i] = v.(*billingpb.RoyaltyReport)
	}

	return objs, nil
}

func (r *royaltyReportRepository) FindCountByMerchantStatusDates(
	ctx context.Context, merchantId string, status []string, dateFrom, dateTo int64,
) (int64, error) {
	var err error
	query := bson.M{}

	if merchantId != "" {
		query["merchant_id"], err = primitive.ObjectIDFromHex(merchantId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return int64(0), err
		}
	}

	if len(status) > 0 {
		query["status"] = bson.M{"$in": status}
	}

	if dateFrom > 0 || dateTo > 0 {
		date := bson.M{}
		if dateFrom > 0 {
			date["$gte"] = time.Unix(dateFrom, 0)
		}
		if dateTo > 0 {
			date["$lte"] = time.Unix(dateTo, 0)
		}
		query["created_at"] = date
	}

	count, err := r.db.Collection(CollectionRoyaltyReport).CountDocuments(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return int64(0), err
	}

	return count, nil
}

func (r *royaltyReportRepository) GetBalanceAmount(ctx context.Context, merchantId, currency string) (float64, error) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return float64(0), err
	}

	royaltyReportsStatusForBalance := []string{
		billingpb.RoyaltyReportStatusAccepted,
		billingpb.RoyaltyReportStatusWaitForPayment,
		billingpb.RoyaltyReportStatusPaid,
	}
	query := []bson.M{
		{
			"$match": bson.M{
				"merchant_id": oid,
				"currency":    currency,
				"status":      bson.M{"$in": royaltyReportsStatusForBalance},
			},
		},
		{
			"$group": bson.M{
				"_id":               "currency",
				"payout_amount":     bson.M{"$sum": "$totals.payout_amount"},
				"correction_amount": bson.M{"$sum": "$totals.correction_amount"},
			},
		},
		{
			"$project": bson.M{
				"_id":    0,
				"amount": bson.M{"$subtract": []interface{}{"$payout_amount", "$correction_amount"}},
			},
		},
	}

	cursor, err := r.db.Collection(CollectionRoyaltyReport).Aggregate(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return 0, err
	}

	res := &pkg2.BalanceQueryResItem{}

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			zap.L().Error(
				pkg.ErrorQueryCursorCloseFailed,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			)
		}
	}()

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			zap.L().Error(
				pkg.ErrorQueryCursorExecutionFailed,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
				zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			)
			return 0, err
		}
	}

	return res.Amount, nil
}

func (r *royaltyReportRepository) GetReportExists(
	ctx context.Context,
	merchantId, currency string,
	from, to time.Time,
) (report *billingpb.RoyaltyReport) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil
	}

	var mgo = models.MgoRoyaltyReport{}
	query := bson.M{
		"merchant_id": oid,
		"period_from": bson.M{"$gte": from},
		"period_to":   bson.M{"$lte": to},
		"currency":    currency,
	}
	err = r.db.Collection(CollectionRoyaltyReport).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String("collection", CollectionRoyaltyReport),
			zap.Any("query", query),
		)
		return nil
	}

	obj, err := r.mapper.MapMgoToObject(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil
	}

	return obj.(*billingpb.RoyaltyReport)
}

func (r *royaltyReportRepository) Insert(ctx context.Context, rr *billingpb.RoyaltyReport, ip, source string) (err error) {
	mgo, err := r.mapper.MapObjectToMgo(rr)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, rr),
		)
		return err
	}

	_, err = r.db.Collection(CollectionRoyaltyReport).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldDocument, rr),
		)
		return
	}

	if err = r.onRoyaltyReportChange(ctx, rr, ip, source); err != nil {
		return
	}

	key := fmt.Sprintf(cacheKeyRoyaltyReport, rr.Id)

	if err = r.cache.Set(key, rr, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, rr),
		)
	}

	return
}

func (r *royaltyReportRepository) Update(ctx context.Context, rr *billingpb.RoyaltyReport, ip, source string) error {
	oid, err := primitive.ObjectIDFromHex(rr.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, rr.Id),
		)
		return nil
	}

	mgo, err := r.mapper.MapObjectToMgo(rr)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, rr),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(CollectionRoyaltyReport).ReplaceOne(ctx, filter, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldDocument, mgo),
		)

		return err
	}

	if err = r.onRoyaltyReportChange(ctx, rr, ip, source); err != nil {
		return err
	}

	key := fmt.Sprintf(cacheKeyRoyaltyReport, rr.Id)

	if err = r.cache.Set(key, rr, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, rr),
		)
	}

	return nil
}

func (r *royaltyReportRepository) UpdateMany(ctx context.Context, query bson.M, set bson.M) error {
	_, err := r.db.Collection(CollectionRoyaltyReport).UpdateMany(context.TODO(), query, set)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSet, set),
		)

		return err
	}

	return nil
}

func (r *royaltyReportRepository) GetById(ctx context.Context, id string) (rr *billingpb.RoyaltyReport, err error) {
	var c = billingpb.RoyaltyReport{}
	key := fmt.Sprintf(cacheKeyRoyaltyReport, id)

	if err := r.cache.Get(key, &c); err == nil {
		return &c, nil
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return
	}

	var mgo = models.MgoRoyaltyReport{}
	filter := bson.M{"_id": oid}
	err = r.db.Collection(CollectionRoyaltyReport).FindOne(ctx, filter).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReport),
			zap.String(pkg.ErrorDatabaseFieldDocumentId, id),
		)
		return
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

	rr = obj.(*billingpb.RoyaltyReport)

	if err = r.cache.Set(key, rr, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, rr),
		)
		// suppress error returning here
		err = nil
	}

	return
}

func (r *royaltyReportRepository) onRoyaltyReportChange(
	ctx context.Context,
	document *billingpb.RoyaltyReport,
	ip, source string,
) (err error) {
	change := &billingpb.RoyaltyReportChanges{
		Id:              primitive.NewObjectID().Hex(),
		RoyaltyReportId: document.Id,
		Source:          source,
		Ip:              ip,
	}

	b, err := json.Marshal(document)

	if err != nil {
		zap.L().Error(
			pkg.ErrorJsonMarshallingFailed,
			zap.Error(err),
			zap.Any("document", document),
		)
		return
	}
	hash := md5.New()
	hash.Write(b)
	change.Hash = hex.EncodeToString(hash.Sum(nil))

	mgo, err := models.NewRoyaltyReportChangesMapper().MapObjectToMgo(change)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, change),
		)
		return err
	}

	_, err = r.db.Collection(CollectionRoyaltyReportChanges).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionRoyaltyReportChanges),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldDocument, mgo),
		)
		return
	}

	return
}
