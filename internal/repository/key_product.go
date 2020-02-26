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
	collectionKeyProduct = "key_product"
)

type keyProductRepository repository

// NewKeyProductRepository create and return an object for working with the key product repository.
// The returned object implements the KeyProductRepositoryInterface interface.
func NewKeyProductRepository(db mongodb.SourceInterface) KeyProductRepositoryInterface {
	s := &keyProductRepository{db: db, mapper: models.NewKeyProductMapper()}
	return s
}

func (r *keyProductRepository) Upsert(ctx context.Context, keyProduct *billingpb.KeyProduct) error {
	oid, err := primitive.ObjectIDFromHex(keyProduct.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, keyProduct.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(keyProduct)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, keyProduct),
		)
		return err
	}

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(collectionKeyProduct).ReplaceOne(ctx, filter, mgo, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpsert),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, keyProduct),
		)
		return err
	}

	return nil
}

func (r *keyProductRepository) Update(ctx context.Context, keyProduct *billingpb.KeyProduct) error {
	oid, _ := primitive.ObjectIDFromHex(keyProduct.Id)

	mgo, err := r.mapper.MapObjectToMgo(keyProduct)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, keyProduct),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(collectionKeyProduct).ReplaceOne(ctx, filter, mgo)

	return err
}

func (r *keyProductRepository) GetById(ctx context.Context, id string) (*billingpb.KeyProduct, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	query := bson.M{"_id": oid, "deleted": false}

	var mgo = models.MgoKeyProduct{}
	err = r.db.Collection(collectionKeyProduct).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, mgo),
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

	return obj.(*billingpb.KeyProduct), nil
}

func (r *keyProductRepository) CountByProjectIdSku(ctx context.Context, projectId, sku string) (int64, error) {
	projectOid, err := primitive.ObjectIDFromHex(projectId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
		)
		return int64(0), err
	}

	query := bson.M{
		"project_id": projectOid,
		"sku":        sku,
		"deleted":    false,
	}
	count, err := r.db.Collection(collectionKeyProduct).CountDocuments(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return int64(0), err
	}

	return count, nil
}

func (r *keyProductRepository) FindByIdsProjectId(ctx context.Context, ids []string, projectId string) ([]*billingpb.KeyProduct, error) {
	var items []primitive.ObjectID

	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
				zap.String(pkg.ErrorDatabaseFieldQuery, id),
			)
			continue
		}

		items = append(items, oid)
	}

	projectOid, err := primitive.ObjectIDFromHex(projectId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
		)
		return nil, err
	}

	query := bson.M{
		"_id":        bson.M{"$in": items},
		"enabled":    true,
		"deleted":    false,
		"project_id": projectOid,
	}
	cursor, err := r.db.Collection(collectionKeyProduct).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoKeyProduct
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.KeyProduct, len(list))

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
		objs[i] = v.(*billingpb.KeyProduct)
	}

	return objs, nil
}

func (r *keyProductRepository) Find(
	ctx context.Context, merchantId, projectId, sku, name, enabled string, offset, limit int64,
) ([]*billingpb.KeyProduct, error) {
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	query := bson.M{"merchant_id": merchantOid, "deleted": false}

	if projectId != "" {
		query["project_id"], err = primitive.ObjectIDFromHex(projectId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
				zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
			)
			return nil, err
		}
	}

	if sku != "" {
		query["sku"] = primitive.Regex{Pattern: sku, Options: "i"}
	}
	if name != "" {
		query["name"] = bson.M{"$elemMatch": bson.M{"value": primitive.Regex{Pattern: name, Options: "i"}}}
	}

	if enabled == "true" {
		query["enabled"] = bson.M{"$eq": true}
	} else if enabled == "false" {
		query["enabled"] = bson.M{"$eq": false}
	}

	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit)
	cursor, err := r.db.Collection(collectionKeyProduct).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoKeyProduct
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.KeyProduct, len(list))

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
		objs[i] = v.(*billingpb.KeyProduct)
	}

	return objs, nil
}

func (r *keyProductRepository) FindCount(ctx context.Context, merchantId, projectId, sku, name, enabled string) (int64, error) {
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return int64(0), err
	}

	query := bson.M{"merchant_id": merchantOid, "deleted": false}

	if projectId != "" {
		query["project_id"], err = primitive.ObjectIDFromHex(projectId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
				zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
			)
			return int64(0), err
		}
	}

	if sku != "" {
		query["sku"] = primitive.Regex{Pattern: sku, Options: "i"}
	}
	if name != "" {
		query["name"] = bson.M{"$elemMatch": bson.M{"value": primitive.Regex{Pattern: name, Options: "i"}}}
	}

	if enabled == "true" {
		query["enabled"] = bson.M{"$eq": true}
	} else if enabled == "false" {
		query["enabled"] = bson.M{"$eq": false}
	}

	count, err := r.db.Collection(collectionKeyProduct).CountDocuments(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionKeyProduct),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return int64(0), err
	}

	return count, nil
}
