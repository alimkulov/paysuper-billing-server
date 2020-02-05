package repository

import (
	"context"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionUserProfile = "user_profile"
)

type userProfileRepository repository

// NewUserProfileRepository create and return an object for working with the user profile repository.
// The returned object implements the UserProfileRepositoryInterface interface.
func NewUserProfileRepository(db mongodb.SourceInterface) UserProfileRepositoryInterface {
	s := &userProfileRepository{db: db}
	return s
}

func (r *userProfileRepository) Add(ctx context.Context, profile *billingpb.UserProfile) error {
	_, err := r.db.Collection(collectionUserProfile).InsertOne(ctx, profile)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, profile),
		)
		return err
	}

	return nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *billingpb.UserProfile) error {
	oid, err := primitive.ObjectIDFromHex(profile.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.String(pkg.ErrorDatabaseFieldQuery, profile.Id),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	err = r.db.Collection(collectionUserProfile).FindOneAndReplace(ctx, filter, profile).Err()

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.String(pkg.ErrorDatabaseFieldQuery, profile.Id),
		)
		return err
	}

	return nil
}

func (r *userProfileRepository) Upsert(ctx context.Context, profile *billingpb.UserProfile) error {
	oid, err := primitive.ObjectIDFromHex(profile.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.String(pkg.ErrorDatabaseFieldQuery, profile.Id),
		)
		return err
	}

	filter := bson.M{"_id": oid}
	opts := options.Replace().SetUpsert(true)
	_, err = r.db.Collection(collectionUserProfile).ReplaceOne(ctx, filter, profile, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.Any(pkg.ErrorDatabaseFieldQuery, profile),
		)

		return err
	}

	return nil
}

func (r *userProfileRepository) GetById(ctx context.Context, id string) (*billingpb.UserProfile, error) {
	var c *billingpb.UserProfile
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.String(pkg.ErrorDatabaseFieldQuery, id),
		)
		return nil, err
	}

	filter := bson.M{"_id": oid}
	err = r.db.Collection(collectionUserProfile).FindOne(ctx, filter).Decode(&c)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.Any(pkg.ErrorDatabaseFieldQuery, filter),
		)
		return nil, err
	}

	return c, nil
}

func (r *userProfileRepository) GetByUserId(ctx context.Context, userId string) (*billingpb.UserProfile, error) {
	var c *billingpb.UserProfile

	filter := bson.M{"user_id": userId}
	err := r.db.Collection(collectionUserProfile).FindOne(ctx, filter).Decode(&c)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionUserProfile),
			zap.Any(pkg.ErrorDatabaseFieldQuery, filter),
		)
		return nil, err
	}

	return c, nil
}
