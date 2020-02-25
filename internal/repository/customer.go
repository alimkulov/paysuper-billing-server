package repository

import (
	"context"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionCustomer = "customer"
)

type customerRepository repository

// NewCustomerRepository create and return an object for working with the customer repository.
// The returned object implements the CustomerRepositoryInterface interface.
func NewCustomerRepository(db mongodb.SourceInterface) CustomerRepositoryInterface {
	s := &customerRepository{db: db, mapper: models.NewCustomerMapper()}
	return s
}

func (r *customerRepository) Insert(ctx context.Context, obj *billingpb.Customer) error {
	mgo, err := r.mapper.MapObjectToMgo(obj)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	_, err = r.db.Collection(collectionCustomer).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	return nil
}

func (r customerRepository) Update(ctx context.Context, obj *billingpb.Customer) error {
	oid, err := primitive.ObjectIDFromHex(obj.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
			zap.String(pkg.ErrorDatabaseFieldQuery, obj.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(obj)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(collectionCustomer).ReplaceOne(ctx, filter, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return err
	}

	return nil
}

func (r customerRepository) GetById(ctx context.Context, id string) (*billingpb.Customer, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	var mgo = models.MgoPayoutDocument{}
	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionCustomer).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
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

	return obj.(*billingpb.Customer), nil
}

func (r customerRepository) Find(ctx context.Context, merchantId string, user *billingpb.TokenUser) (*billingpb.Customer, error) {
	var (
		err         error
		merchantOid primitive.ObjectID
	)

	if merchantId != "" {
		merchantOid, err = primitive.ObjectIDFromHex(merchantId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return nil, err
		}
	}

	var subQuery []bson.M
	var subQueryItem bson.M

	if user.Id != "" {
		subQueryItem = bson.M{
			"identity": bson.M{
				"$elemMatch": bson.M{
					"type":        pkg.UserIdentityTypeExternal,
					"merchant_id": merchantOid,
					"value":       user.Id,
				},
			},
		}
		subQuery = append(subQuery, subQueryItem)
	}

	if user.Email != nil && user.Email.Value != "" {
		subQueryItem = bson.M{
			"identity": bson.M{
				"$elemMatch": bson.M{
					"type":        pkg.UserIdentityTypeEmail,
					"merchant_id": merchantOid,
					"value":       user.Email.Value,
				},
			},
		}
		subQuery = append(subQuery, subQueryItem)
	}

	if user.Phone != nil && user.Phone.Value != "" {
		subQueryItem = bson.M{
			"identity": bson.M{
				"$elemMatch": bson.M{
					"type":        pkg.UserIdentityTypePhone,
					"merchant_id": merchantOid,
					"value":       user.Phone.Value,
				},
			},
		}
		subQuery = append(subQuery, subQueryItem)
	}

	query := make(bson.M)

	if subQuery == nil || len(subQuery) <= 0 {
		return nil, mongo.ErrNoDocuments
	}

	if len(subQuery) > 1 {
		query["$or"] = subQuery
	} else {
		query = subQuery[0]
	}

	mgo := models.MgoCustomer{}
	err = r.db.Collection(collectionCustomer).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionCustomer),
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

	return obj.(*billingpb.Customer), nil
}
