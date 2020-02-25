package repository

import (
	"context"
	"fmt"
	"github.com/paysuper/paysuper-billing-server/internal/database"
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
	collectionOperatingCompanies = "operating_companies"

	cacheKeyOperatingCompany                 = "operating_company:id:%s"
	cacheKeyOperatingCompanyByPaymentCountry = "operating_company:country:%s"
	cacheKeyAllOperatingCompanies            = "operating_company:all"
)

type operatingCompanyRepository repository

// NewOperatingCompanyRepository create and return an object for working with the operating company repository.
// The returned object implements the OperatingCompanyRepositoryInterface interface.
func NewOperatingCompanyRepository(db mongodb.SourceInterface, cache database.CacheInterface) OperatingCompanyRepositoryInterface {
	s := &operatingCompanyRepository{db: db, cache: cache, mapper: models.NewOperatingCompanyMapper()}
	return s
}

func (r *operatingCompanyRepository) GetById(ctx context.Context, id string) (*billingpb.OperatingCompany, error) {
	oc := &billingpb.OperatingCompany{}
	key := fmt.Sprintf(cacheKeyOperatingCompany, id)

	if err := r.cache.Get(key, oc); err == nil {
		return oc, nil
	}

	oid, _ := primitive.ObjectIDFromHex(id)

	mgo := &models.MgoOperatingCompany{}
	filter := bson.M{"_id": oid}
	err := r.db.Collection(collectionOperatingCompanies).FindOne(ctx, filter).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
			zap.String(pkg.ErrorDatabaseFieldDocumentId, id),
		)
		return nil, err
	}

	obj, err := r.mapper.MapMgoToObject(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	oc = obj.(*billingpb.OperatingCompany)

	err = r.cache.Set(key, oc, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, oc),
		)
	}

	return oc, nil
}

func (r *operatingCompanyRepository) GetByPaymentCountry(ctx context.Context, code string) (*billingpb.OperatingCompany, error) {
	oc := &billingpb.OperatingCompany{}
	key := fmt.Sprintf(cacheKeyOperatingCompanyByPaymentCountry, code)

	if err := r.cache.Get(key, oc); err == nil {
		return oc, nil
	}

	query := bson.M{"payment_countries": code}

	if code == "" {
		query["payment_countries"] = bson.M{"$size": 0}
	} else {
		query["payment_countries"] = code
	}

	mgo := &models.MgoOperatingCompany{}
	err := r.db.Collection(collectionOperatingCompanies).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := r.mapper.MapMgoToObject(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
		)
		return nil, err
	}

	oc = obj.(*billingpb.OperatingCompany)

	err = r.cache.Set(key, oc, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, oc),
		)
	}

	return oc, err
}

func (r *operatingCompanyRepository) GetAll(ctx context.Context) ([]*billingpb.OperatingCompany, error) {
	var result []*billingpb.OperatingCompany
	err := r.cache.Get(cacheKeyAllOperatingCompanies, &result)

	if err == nil {
		return result, nil
	}

	cursor, err := r.db.Collection(collectionOperatingCompanies).Find(ctx, bson.M{})

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
		)
		return nil, err
	}

	var list []*models.MgoOperatingCompany
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
		)
		return nil, err
	}

	objs := make([]*billingpb.OperatingCompany, len(list))

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
		objs[i] = v.(*billingpb.OperatingCompany)
	}

	err = r.cache.Set(cacheKeyAllOperatingCompanies, objs, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, cacheKeyAllOperatingCompanies),
			zap.Any(pkg.ErrorCacheFieldData, objs),
		)
	}

	return objs, nil
}

func (r *operatingCompanyRepository) Upsert(ctx context.Context, oc *billingpb.OperatingCompany) error {
	oid, err := primitive.ObjectIDFromHex(oc.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
			zap.String(pkg.ErrorDatabaseFieldQuery, oc.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(oc)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, oc),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	opts := options.Replace().SetUpsert(true)
	_, err = r.db.Collection(collectionOperatingCompanies).ReplaceOne(ctx, filter, mgo, opts)

	if err != nil {
		zap.S().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionOperatingCompanies),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpsert),
			zap.Any(pkg.ErrorDatabaseFieldDocument, oc),
		)
		return err
	}

	_ = r.updateCache(oc)

	return nil
}

func (r *operatingCompanyRepository) Exists(ctx context.Context, id string) bool {
	oc, err := r.GetById(ctx, id)

	return err == nil && oc != nil
}

func (r operatingCompanyRepository) updateCache(oc *billingpb.OperatingCompany) error {
	key := fmt.Sprintf(cacheKeyOperatingCompany, oc.Id)
	err := r.cache.Set(key, oc, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, oc),
		)
		return err
	}

	if len(oc.Country) == 0 {
		key = fmt.Sprintf(cacheKeyOperatingCompanyByPaymentCountry, "")
		err = r.cache.Set(key, oc, 0)

		if err != nil {
			zap.L().Error(
				pkg.ErrorCacheQueryFailed,
				zap.Error(err),
				zap.String(pkg.ErrorCacheFieldCmd, "SET"),
				zap.String(pkg.ErrorCacheFieldKey, key),
				zap.Any(pkg.ErrorCacheFieldData, oc),
			)
			return err
		}
	} else {
		for _, countryCode := range oc.PaymentCountries {
			key = fmt.Sprintf(cacheKeyOperatingCompanyByPaymentCountry, countryCode)
			err = r.cache.Set(key, oc, 0)

			if err != nil {
				zap.L().Error(
					pkg.ErrorCacheQueryFailed,
					zap.Error(err),
					zap.String(pkg.ErrorCacheFieldCmd, "SET"),
					zap.String(pkg.ErrorCacheFieldKey, key),
					zap.Any(pkg.ErrorCacheFieldData, oc),
				)
				return err
			}
		}
	}

	err = r.cache.Delete(cacheKeyAllOperatingCompanies)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "DELETE"),
			zap.String(pkg.ErrorCacheFieldKey, cacheKeyAllOperatingCompanies),
		)
	}

	return nil
}
