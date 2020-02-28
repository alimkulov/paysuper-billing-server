package repository

import (
	"context"
	"fmt"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	"github.com/paysuper/paysuper-billing-server/internal/helper"
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
	// CollectionCountry is name of table for collection the country.
	CollectionCountry = "country"

	cacheCountryCodeA2           = "country:code_a2:%s"
	cacheCountryRisk             = "country:risk:%t"
	cacheCountryAll              = "country:all"
	cacheCountryRegions          = "country:regions"
	cacheCountriesWithVatEnabled = "country:with_vat"
)

type countryRepository repository

// NewCountryRepository create and return an object for working with the country repository.
// The returned object implements the CountryRepositoryInterface interface.
func NewCountryRepository(db mongodb.SourceInterface, cache database.CacheInterface) CountryRepositoryInterface {
	s := &countryRepository{db: db, cache: cache, mapper: models.NewCountryMapper()}
	return s
}

func (h *countryRepository) Insert(ctx context.Context, country *billingpb.Country) error {
	mgo, err := h.mapper.MapObjectToMgo(country)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	_, err = h.db.Collection(CollectionCountry).InsertOne(ctx, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	if err := h.updateCache(country); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	return nil
}

func (h *countryRepository) MultipleInsert(ctx context.Context, country []*billingpb.Country) error {
	c := make([]interface{}, len(country))
	for i, v := range country {
		if v.VatThreshold == nil {
			v.VatThreshold = &billingpb.CountryVatThreshold{
				Year:  0,
				World: 0,
			}
		}
		v.VatThreshold.Year = tools.FormatAmount(v.VatThreshold.Year)
		v.VatThreshold.World = tools.FormatAmount(v.VatThreshold.World)
		mgo, err := h.mapper.MapObjectToMgo(v)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, v),
			)
			return err
		}
		c[i] = mgo
	}

	_, err := h.db.Collection(CollectionCountry).InsertMany(ctx, c)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldQuery, c),
		)
		return err
	}

	if err := h.updateCache(nil); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	return nil
}

func (h *countryRepository) Update(ctx context.Context, country *billingpb.Country) error {
	country.VatThreshold.Year = tools.FormatAmount(country.VatThreshold.Year)
	country.VatThreshold.World = tools.FormatAmount(country.VatThreshold.World)

	oid, err := primitive.ObjectIDFromHex(country.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.String(pkg.ErrorDatabaseFieldQuery, country.Id),
		)
		return err
	}

	mgo, err := h.mapper.MapObjectToMgo(country)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	_, err = h.db.Collection(CollectionCountry).ReplaceOne(ctx, bson.M{"_id": oid}, mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	if err := h.updateCache(country); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, country),
		)
		return err
	}

	return nil
}

func (h *countryRepository) GetByIsoCodeA2(ctx context.Context, code string) (*billingpb.Country, error) {
	c := &billingpb.Country{}
	var key = fmt.Sprintf(cacheCountryCodeA2, code)

	if err := h.cache.Get(key, c); err == nil {
		return c, nil
	}


	mgo := &models.MgoCountry{}
	query := bson.M{"iso_code_a2": code}
	err := h.db.Collection(CollectionCountry).FindOne(ctx, query).Decode(mgo)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	obj, err := h.mapper.MapMgoToObject(mgo)
	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldQuery, obj),
		)
		return nil, err
	}
	c = obj.(*billingpb.Country)

	if err = h.cache.Set(key, c, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorDatabaseFieldQuery, c),
		)
	}

	return c, nil
}

func (h *countryRepository) FindByVatEnabled(ctx context.Context) (*billingpb.CountriesList, error) {
	var c = &billingpb.CountriesList{}
	key := cacheCountriesWithVatEnabled

	if err := h.cache.Get(key, c); err == nil {
		return c, nil
	}

	query := bson.M{
		"vat_enabled":               true,
		"iso_code_a2":               bson.M{"$ne": "US"},
		"vat_currency_rates_policy": bson.M{"$ne": ""},
		"vat_period_month":          bson.M{"$gt": 0},
	}
	opts := options.Find().SetSort(bson.M{"iso_code_a2": 1})
	cursor, err := h.db.Collection(CollectionCountry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var mgoCountries []*models.MgoCountry
	err = cursor.All(ctx, &mgoCountries)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	c.Countries = make([]*billingpb.Country, len(mgoCountries))
	for i, country := range mgoCountries {
		obj, err := h.mapper.MapMgoToObject(country)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, country),
			)
			return nil, err
		}
		c.Countries[i] = obj.(*billingpb.Country)
	}

	if err = h.cache.Set(key, c, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, c),
		)
	}

	return c, nil
}

func (h *countryRepository) GetAll(ctx context.Context) (*billingpb.CountriesList, error) {
	var c = &billingpb.CountriesList{}
	key := cacheCountryAll

	if err := h.cache.Get(key, c); err == nil {
		return c, nil
	}

	query := bson.M{}
	opts := options.Find().SetSort(bson.M{"iso_code_a2": 1})
	cursor, err := h.db.Collection(CollectionCountry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var mgoCountries []*models.MgoCountry
	err = cursor.All(ctx, &mgoCountries)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	c.Countries = make([]*billingpb.Country, len(mgoCountries))
	for i, country := range mgoCountries {
		obj, err := h.mapper.MapMgoToObject(country)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, country),
			)
			return nil, err
		}
		c.Countries[i] = obj.(*billingpb.Country)
	}


	if err = h.cache.Set(key, c, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorCacheFieldData, c),
		)
	}

	return c, nil
}

func (h *countryRepository) FindByHighRisk(ctx context.Context, isHighRiskOrder bool) (*billingpb.CountriesList, error) {
	var c = &billingpb.CountriesList{}
	var key = fmt.Sprintf(cacheCountryRisk, isHighRiskOrder)

	if err := h.cache.Get(key, c); err == nil {
		return c, nil
	}

	field := "payments_allowed"
	if isHighRiskOrder {
		field = "high_risk_payments_allowed"
	}

	query := bson.M{
		field: true,
	}
	opts := options.Find().SetSort(bson.M{"iso_code_a2": 1})
	cursor, err := h.db.Collection(CollectionCountry).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	var mgoCountries []*models.MgoCountry
	err = cursor.All(ctx, &mgoCountries)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, CollectionCountry),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return nil, err
	}

	c.Countries = make([]*billingpb.Country, len(mgoCountries))
	for i, country := range mgoCountries {
		obj, err := h.mapper.MapMgoToObject(country)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, country),
			)
			return nil, err
		}
		c.Countries[i] = obj.(*billingpb.Country)
	}

	if err = h.cache.Set(key, c, 0); err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key),
			zap.Any(pkg.ErrorDatabaseFieldQuery, c),
		)
	}

	return c, nil
}

func (h *countryRepository) IsTariffRegionSupported(region string) bool {
	return helper.Contains(pkg.SupportedTariffRegions, region)
}

func (h *countryRepository) updateCache(country *billingpb.Country) error {
	if country != nil {
		if err := h.cache.Set(fmt.Sprintf(cacheCountryCodeA2, country.IsoCodeA2), country, 0); err != nil {
			return err
		}
	}

	if err := h.cache.Delete(cacheCountryAll); err != nil {
		return err
	}

	if err := h.cache.Delete(cacheCountryRegions); err != nil {
		return err
	}

	if err := h.cache.Delete(cacheCountriesWithVatEnabled); err != nil {
		return err
	}

	if err := h.cache.Delete(fmt.Sprintf(cacheCountryRisk, false)); err != nil {
		return err
	}

	if err := h.cache.Delete(fmt.Sprintf(cacheCountryRisk, true)); err != nil {
		return err
	}

	return nil
}
