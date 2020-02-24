package repository

import (
	"context"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionNotification = "notification"
)

type notificationRepository repository

// NewNotificationRepository create and return an object for working with the notification repository.
// The returned object implements the NotificationRepositoryInterface interface.
func NewNotificationRepository(db mongodb.SourceInterface) NotificationRepositoryInterface {
	s := &notificationRepository{db: db, mapper: models.NewNotificationMapper()}
	return s
}

func (r *notificationRepository) Insert(ctx context.Context, obj *billingpb.Notification) error {
	mgo, err := r.mapper.MapObjectToMgo(obj)

	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	_, err = r.db.Collection(collectionNotification).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return err
	}

	return nil
}

func (r *notificationRepository) Update(ctx context.Context, obj *billingpb.Notification) error {
	oid, err := primitive.ObjectIDFromHex(obj.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.String(pkg.ErrorDatabaseFieldQuery, obj.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(obj)

	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(collectionNotification).ReplaceOne(ctx, filter, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return err
	}

	return nil
}

func (r notificationRepository) GetById(ctx context.Context, id string) (*billingpb.Notification, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	mgo := &models.MgoNotification{}
	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionNotification).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	decoded, err := r.mapper.MapMgoToObject(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	return decoded.(*billingpb.Notification), nil
}

func (r *notificationRepository) Find(
	ctx context.Context,
	merchantId,
	userId string,
	isSystem int32,
	sort []string,
	offset,
	limit int64,
) ([]*billingpb.Notification, error) {
	var err error

	query := make(bson.M)

	if merchantId != "" {
		query["merchant_id"], err = primitive.ObjectIDFromHex(merchantId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return nil, err
		}
	}

	if userId != "" {
		query["user_id"] = userId
	}

	if isSystem > 0 {
		if isSystem == 1 {
			query["is_system"] = false
		} else {
			query["is_system"] = true
		}
	}

	opts := options.Find().
		SetSort(mongodb.ToSortOption(sort)).
		SetLimit(limit).
		SetSkip(offset)
	cursor, err := r.db.Collection(collectionNotification).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var mgo []*models.MgoNotification

	err = cursor.All(ctx, &mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionProject),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	if len(mgo) == 0 {
		return nil, nil
	}

	notifications := make([]*billingpb.Notification, len(mgo))
	for i, mgoNotif := range mgo {
		obj, err := r.mapper.MapMgoToObject(mgoNotif)
		if err != nil {
			zap.L().Error(
				pkg.ErrorMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, mgoNotif),
			)
			return nil, err
		}

		notifications[i] = obj.(*billingpb.Notification)
	}

	return notifications, nil
}

func (r *notificationRepository) FindCount(ctx context.Context, merchantId, userId string, isSystem int32) (int64, error) {
	var err error

	query := make(bson.M)

	if merchantId != "" {
		query["merchant_id"], err = primitive.ObjectIDFromHex(merchantId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return int64(0), err
		}
	}

	if userId != "" {
		query["user_id"] = userId
	}

	if isSystem > 0 {
		if isSystem == 1 {
			query["is_system"] = false
		} else {
			query["is_system"] = true
		}
	}

	count, err := r.db.Collection(collectionNotification).CountDocuments(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionNotification),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return int64(0), err
	}

	return count, nil
}
