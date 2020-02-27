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
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionMerchantUsersTable = "user_merchant"
	collectionAdminUsersTable    = "user_admin"

	cacheUserMerchants = "user:merchants:%s"
)

type userRoleRepository repository

// NewUserRoleRepository create and return an object for working with the user role repository.
// The returned object implements the UserRoleRepositoryInterface interface.
func NewUserRoleRepository(db mongodb.SourceInterface, cache database.CacheInterface) UserRoleRepositoryInterface {
	s := &userRoleRepository{db: db, cache: cache, mapper: models.NewUserRoleMapper()}
	return s
}

func (r *userRoleRepository) AddMerchantUser(ctx context.Context, role *billingpb.UserRole) error {
	mgo, err := r.mapper.MapObjectToMgo(role)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	_, err = r.db.Collection(collectionMerchantUsersTable).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) AddAdminUser(ctx context.Context, role *billingpb.UserRole) error {
	mgo, err := r.mapper.MapObjectToMgo(role)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	_, err = r.db.Collection(collectionAdminUsersTable).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) UpdateMerchantUser(ctx context.Context, role *billingpb.UserRole) error {
	oid, err := primitive.ObjectIDFromHex(role.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, role.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(role)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	err = r.db.Collection(collectionMerchantUsersTable).FindOneAndReplace(ctx, filter, mgo).Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	key := fmt.Sprintf(cacheUserMerchants, role.UserId)
	if err := r.cache.Delete(key); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorCacheFieldCmd, "delete"),
			zap.Any(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, role.UserId),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) UpdateAdminUser(ctx context.Context, role *billingpb.UserRole) error {
	oid, err := primitive.ObjectIDFromHex(role.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, role.Id),
		)
		return err
	}

	mgo, err := r.mapper.MapObjectToMgo(role)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	err = r.db.Collection(collectionAdminUsersTable).FindOneAndReplace(ctx, filter, mgo).Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, role),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) GetAdminUserById(ctx context.Context, id string) (*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	var mgo = models.MgoUserRole{}
	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionAdminUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetMerchantUserById(ctx context.Context, id string) (*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	var mgo = models.MgoUserRole{}
	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionMerchantUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) DeleteAdminUser(ctx context.Context, role *billingpb.UserRole) error {
	oid, err := primitive.ObjectIDFromHex(role.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, role.Id),
		)
		return err
	}

	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionAdminUsersTable).FindOneAndDelete(ctx, query).Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) DeleteMerchantUser(ctx context.Context, role *billingpb.UserRole) error {
	oid, err := primitive.ObjectIDFromHex(role.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, role.Id),
		)
		return err
	}

	query := bson.M{"_id": oid}
	err = r.db.Collection(collectionMerchantUsersTable).FindOneAndDelete(ctx, query).Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return err
	}

	key := fmt.Sprintf(cacheUserMerchants, role.UserId)
	if err := r.cache.Delete(key); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorCacheFieldCmd, "delete"),
			zap.Any(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, role.UserId),
		)
		return err
	}

	return nil
}

func (r *userRoleRepository) GetSystemAdmin(ctx context.Context) (*billingpb.UserRole, error) {
	var mgo = models.MgoUserRole{}

	query := bson.M{"role": billingpb.RoleSystemAdmin}
	err := r.db.Collection(collectionAdminUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetMerchantOwner(ctx context.Context, merchantId string) (*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	var mgo = models.MgoUserRole{}
	query := bson.M{"merchant_id": oid, "role": billingpb.RoleMerchantOwner}
	err = r.db.Collection(collectionMerchantUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetMerchantsForUser(ctx context.Context, userId string) ([]*billingpb.UserRole, error) {
	var users []*billingpb.UserRole

	if err := r.cache.Get(fmt.Sprintf(cacheUserMerchants, userId), users); err == nil {
		return users, nil
	}

	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, userId),
		)
		return nil, err
	}

	query := bson.M{"user_id": oid}
	cursor, err := r.db.Collection(collectionMerchantUsersTable).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var list []*models.MgoUserRole
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.UserRole, len(list))

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
		objs[i] = v.(*billingpb.UserRole)
	}

	err = r.cache.Set(fmt.Sprintf(cacheUserMerchants, userId), objs, 0)

	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorCacheFieldCmd, "set"),
			zap.Any(pkg.ErrorCacheFieldKey, cacheUserMerchants),
			zap.Any(pkg.ErrorCacheFieldData, objs),
		)
	}

	return objs, nil
}

func (r *userRoleRepository) GetUsersForAdmin(ctx context.Context) ([]*billingpb.UserRole, error) {
	cursor, err := r.db.Collection(collectionAdminUsersTable).Find(ctx, bson.M{})

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
		)

		return nil, err
	}

	var list []*models.MgoUserRole
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
		)
		return nil, err
	}

	objs := make([]*billingpb.UserRole, len(list))

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
		objs[i] = v.(*billingpb.UserRole)
	}

	return objs, nil
}

func (r *userRoleRepository) GetUsersForMerchant(ctx context.Context, merchantId string) ([]*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	query := bson.M{"merchant_id": oid}
	cursor, err := r.db.Collection(collectionMerchantUsersTable).Find(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)

		return nil, err
	}

	var list []*models.MgoUserRole
	err = cursor.All(ctx, &list)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	objs := make([]*billingpb.UserRole, len(list))

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
		objs[i] = v.(*billingpb.UserRole)
	}

	return objs, nil
}

func (r *userRoleRepository) GetAdminUserByEmail(ctx context.Context, email string) (*billingpb.UserRole, error) {
	var mgo = models.MgoUserRole{}

	query := bson.M{"email": email}
	err := r.db.Collection(collectionAdminUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetMerchantUserByEmail(ctx context.Context, merchantId string, email string) (*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	var mgo = models.MgoUserRole{}
	filter := bson.M{"merchant_id": oid, "email": email}
	err = r.db.Collection(collectionMerchantUsersTable).FindOne(ctx, filter).Decode(&mgo)

	if err != nil {
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetAdminUserByUserId(ctx context.Context, userId string) (*billingpb.UserRole, error) {
	var mgo = models.MgoUserRole{}

	oid, _ := primitive.ObjectIDFromHex(userId)
	query := bson.M{"user_id": oid}
	err := r.db.Collection(collectionAdminUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionAdminUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}

func (r *userRoleRepository) GetMerchantUserByUserId(ctx context.Context, merchantId string, id string) (*billingpb.UserRole, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	var mgo = models.MgoUserRole{}
	query := bson.M{"merchant_id": merchantOid, "user_id": oid}
	err = r.db.Collection(collectionMerchantUsersTable).FindOne(ctx, query).Decode(&mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionMerchantUsersTable),
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

	return obj.(*billingpb.UserRole), nil
}
