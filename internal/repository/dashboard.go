package repository

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	pkg2 "github.com/paysuper/paysuper-billing-server/internal/pkg"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"time"
)

const (
	dashboardMainGrossRevenueAndVatCacheKey       = "dashboard:main:gross_revenue_and_vat:%x"
	dashboardMainTotalTransactionsAndArpuCacheKey = "dashboard:main:total_transactions_and_arpu:%x"
	dashboardRevenueDynamicCacheKey               = "dashboard:revenue_dynamic:%x"
	dashboardBaseRevenueByCountryCacheKey         = "dashboard:base:revenue_by_country:%x"
	dashboardBaseSalesTodayCacheKey               = "dashboard:base:sales_today:%x"
	dashboardBaseSourcesCacheKey                  = "dashboard:base:sources:%x"

	dashboardReportGroupByHour        = "$hour"
	dashboardReportGroupByDay         = "$day"
	dashboardReportGroupByMonth       = "$month"
	dashboardReportGroupByWeek        = "$week"
	dashboardReportGroupByPeriodInDay = "$period_in_day"
)

var (
	dashboardReportBasePreviousPeriodsNames = map[string]string{
		pkg.DashboardPeriodCurrentDay:      pkg.DashboardPeriodPreviousDay,
		pkg.DashboardPeriodPreviousDay:     pkg.DashboardPeriodTwoDaysAgo,
		pkg.DashboardPeriodCurrentWeek:     pkg.DashboardPeriodPreviousWeek,
		pkg.DashboardPeriodPreviousWeek:    pkg.DashboardPeriodTwoWeeksAgo,
		pkg.DashboardPeriodCurrentMonth:    pkg.DashboardPeriodPreviousMonth,
		pkg.DashboardPeriodPreviousMonth:   pkg.DashboardPeriodTwoMonthsAgo,
		pkg.DashboardPeriodCurrentQuarter:  pkg.DashboardPeriodPreviousQuarter,
		pkg.DashboardPeriodPreviousQuarter: pkg.DashboardPeriodTwoQuarterAgo,
		pkg.DashboardPeriodCurrentYear:     pkg.DashboardPeriodPreviousYear,
		pkg.DashboardPeriodPreviousYear:    pkg.DashboardPeriodTwoYearsAgo,
	}
)

type dashboardRepository repository

// NewDashboardRepository create and return an object for working with the key repository.
// The returned object implements the DashboardRepositoryInterface interface.
func NewDashboardRepository(db mongodb.SourceInterface, cache database.CacheInterface) DashboardRepositoryInterface {
	s := &dashboardRepository{db: db, cache: cache}
	return s
}

func (r *dashboardRepository) GetMainReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardMainReport, error) {
	processorGrossRevenueAndVatCurrent, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardMainGrossRevenueAndVatCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	dataGrossRevenueAndVatCurrent, err := processorGrossRevenueAndVatCurrent.ExecuteReport(
		ctx,
		new(models.MgoGrossRevenueAndVatReports),
		processorGrossRevenueAndVatCurrent.ExecuteGrossRevenueAndVatReports,
	)

	if err != nil {
		return nil, err
	}

	processorGrossRevenueAndVatPrevious, err := r.newDashboardReportProcessor(
		merchantId,
		dashboardReportBasePreviousPeriodsNames[period],
		dashboardMainGrossRevenueAndVatCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	dataGrossRevenueAndVatPrevious, err := processorGrossRevenueAndVatPrevious.ExecuteReport(
		ctx,
		new(models.MgoGrossRevenueAndVatReports),
		processorGrossRevenueAndVatPrevious.ExecuteGrossRevenueAndVatReports,
	)

	if err != nil {
		return nil, err
	}

	processorTotalTransactionsAndArpuCurrent, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardMainTotalTransactionsAndArpuCacheKey,
		bson.M{"$in": []string{"processed", "refunded", "chargeback"}},
	)

	if err != nil {
		return nil, err
	}

	dataTotalTransactionsAndArpuCurrent, err := processorTotalTransactionsAndArpuCurrent.ExecuteReport(
		ctx,
		new(models.MgoTotalTransactionsAndArpuReports),
		processorTotalTransactionsAndArpuCurrent.ExecuteTotalTransactionsAndArpuReports,
	)

	if err != nil {
		return nil, err
	}

	processorTotalTransactionsAndArpuPrevious, err := r.newDashboardReportProcessor(
		merchantId,
		dashboardReportBasePreviousPeriodsNames[period],
		dashboardMainTotalTransactionsAndArpuCacheKey,
		bson.M{"$in": []string{"processed", "refunded", "chargeback"}},
	)

	if err != nil {
		return nil, err
	}

	dataTotalTransactionsAndArpuPrevious, err := processorTotalTransactionsAndArpuPrevious.ExecuteReport(
		ctx,
		new(models.MgoTotalTransactionsAndArpuReports),
		processorTotalTransactionsAndArpuPrevious.ExecuteTotalTransactionsAndArpuReports,
	)

	if err != nil {
		return nil, err
	}

	dataGrossRevenueAndVatCurrentTyped := dataGrossRevenueAndVatCurrent.(*pkg2.GrossRevenueAndVatReports)
	dataGrossRevenueAndVatPreviousTyped := dataGrossRevenueAndVatPrevious.(*pkg2.GrossRevenueAndVatReports)

	if dataGrossRevenueAndVatPreviousTyped == nil {
		dataGrossRevenueAndVatPreviousTyped = &pkg2.GrossRevenueAndVatReports{
			GrossRevenue: &billingpb.DashboardAmountItemWithChart{},
			Vat:          &billingpb.DashboardAmountItemWithChart{},
		}
	}
	dataGrossRevenueAndVatCurrentTyped.GrossRevenue.AmountPrevious = dataGrossRevenueAndVatPreviousTyped.GrossRevenue.AmountCurrent
	dataGrossRevenueAndVatCurrentTyped.Vat.AmountPrevious = dataGrossRevenueAndVatPreviousTyped.Vat.AmountCurrent

	dataTotalTransactionsAndArpuCurrentTyped := dataTotalTransactionsAndArpuCurrent.(*pkg2.TotalTransactionsAndArpuReports)
	dataTotalTransactionsAndArpuPreviousTyped := dataTotalTransactionsAndArpuPrevious.(*pkg2.TotalTransactionsAndArpuReports)

	if dataTotalTransactionsAndArpuPreviousTyped == nil {
		dataTotalTransactionsAndArpuPreviousTyped = &pkg2.TotalTransactionsAndArpuReports{
			TotalTransactions: &billingpb.DashboardMainReportTotalTransactions{},
			Arpu:              &billingpb.DashboardAmountItemWithChart{},
		}
	}
	dataTotalTransactionsAndArpuCurrentTyped.TotalTransactions.CountPrevious = dataTotalTransactionsAndArpuPreviousTyped.TotalTransactions.CountCurrent
	dataTotalTransactionsAndArpuCurrentTyped.Arpu.AmountPrevious = dataTotalTransactionsAndArpuPreviousTyped.Arpu.AmountCurrent

	result := &billingpb.DashboardMainReport{
		GrossRevenue:      dataGrossRevenueAndVatCurrentTyped.GrossRevenue,
		Vat:               dataGrossRevenueAndVatCurrentTyped.Vat,
		TotalTransactions: dataTotalTransactionsAndArpuCurrentTyped.TotalTransactions,
		Arpu:              dataTotalTransactionsAndArpuCurrentTyped.Arpu,
	}

	return result, nil
}

func (r *dashboardRepository) GetBaseReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardBaseReports, error) {
	revenueByCountryReport, err := r.getBaseRevenueByCountryReport(ctx, merchantId, period)

	if err != nil {
		return nil, err
	}

	salesTodayReport, err := r.getBaseSalesTodayReport(ctx, merchantId, period)

	if err != nil {
		return nil, err
	}

	sourcesReport, err := r.getBaseSourcesReport(ctx, merchantId, period)

	if err != nil {
		return nil, err
	}

	reports := &billingpb.DashboardBaseReports{
		RevenueByCountry: revenueByCountryReport,
		SalesToday:       salesTodayReport,
		Sources:          sourcesReport,
	}

	return reports, nil
}

func (r *dashboardRepository) GetRevenueDynamicsReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardRevenueDynamicReport, error) {
	processor, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardRevenueDynamicCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	data, err := processor.ExecuteReport(
		ctx,
		new(models.MgoDashboardRevenueDynamicReport),
		processor.ExecuteRevenueDynamicReport,
	)

	if err != nil {
		return nil, err
	}

	dataTyped := data.(*billingpb.DashboardRevenueDynamicReport)

	if len(dataTyped.Items) > 0 {
		dataTyped.Currency = dataTyped.Items[0].Currency
	}

	return dataTyped, nil
}

func (r *dashboardRepository) getBaseRevenueByCountryReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardRevenueByCountryReport, error) {
	processorCurrent, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardBaseRevenueByCountryCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	processorPrevious, err := r.newDashboardReportProcessor(
		merchantId,
		dashboardReportBasePreviousPeriodsNames[period],
		dashboardBaseRevenueByCountryCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	dataCurrent, err := processorCurrent.ExecuteReport(
		ctx,
		new(models.MgoDashboardRevenueByCountryReport),
		processorCurrent.ExecuteRevenueByCountryReport,
	)

	if err != nil {
		return nil, err
	}

	dataPrevious, err := processorPrevious.ExecuteReport(
		ctx,
		new(models.MgoDashboardRevenueByCountryReport),
		processorPrevious.ExecuteRevenueByCountryReport,
	)

	if err != nil {
		return nil, err
	}

	dataCurrentTyped := dataCurrent.(*billingpb.DashboardRevenueByCountryReport)

	if dataCurrentTyped == nil {
		dataCurrentTyped = &billingpb.DashboardRevenueByCountryReport{}
	}

	dataPreviousTyped := dataPrevious.(*billingpb.DashboardRevenueByCountryReport)

	if dataPreviousTyped == nil {
		dataPreviousTyped = &billingpb.DashboardRevenueByCountryReport{}
	}

	dataCurrentTyped.TotalPrevious = dataPreviousTyped.TotalCurrent

	return dataCurrentTyped, nil
}

func (r *dashboardRepository) getBaseSalesTodayReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardSalesTodayReport, error) {
	processorCurrent, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardBaseSalesTodayCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	processorPrevious, err := r.newDashboardReportProcessor(
		merchantId,
		dashboardReportBasePreviousPeriodsNames[period],
		dashboardBaseSalesTodayCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	dataCurrent, err := processorCurrent.ExecuteReport(
		ctx,
		new(billingpb.DashboardSalesTodayReport),
		processorCurrent.ExecuteSalesTodayReport,
	)

	if err != nil {
		return nil, err
	}

	dataPrevious, err := processorPrevious.ExecuteReport(
		ctx,
		new(billingpb.DashboardSalesTodayReport),
		processorPrevious.ExecuteSalesTodayReport,
	)

	if err != nil {
		return nil, err
	}

	dataCurrentTyped := dataCurrent.(*billingpb.DashboardSalesTodayReport)
	dataPreviousTyped := dataPrevious.(*billingpb.DashboardSalesTodayReport)
	dataCurrentTyped.TotalPrevious = dataPreviousTyped.TotalCurrent

	return dataCurrentTyped, nil
}

func (r *dashboardRepository) getBaseSourcesReport(
	ctx context.Context,
	merchantId, period string,
) (*billingpb.DashboardSourcesReport, error) {
	processorCurrent, err := r.newDashboardReportProcessor(
		merchantId,
		period,
		dashboardBaseSourcesCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	processorPrevious, err := r.newDashboardReportProcessor(
		merchantId,
		dashboardReportBasePreviousPeriodsNames[period],
		dashboardBaseSourcesCacheKey,
		"processed",
	)

	if err != nil {
		return nil, err
	}

	dataCurrent, err := processorCurrent.ExecuteReport(
		ctx,
		new(billingpb.DashboardSourcesReport),
		processorCurrent.ExecuteSourcesReport,
	)

	if err != nil {
		return nil, err
	}

	dataPrevious, err := processorPrevious.ExecuteReport(
		ctx,
		new(billingpb.DashboardSourcesReport),
		processorPrevious.ExecuteSourcesReport,
	)

	if err != nil {
		return nil, err
	}

	dataCurrentTyped := dataCurrent.(*billingpb.DashboardSourcesReport)
	dataPreviousTyped := dataPrevious.(*billingpb.DashboardSourcesReport)
	dataCurrentTyped.TotalPrevious = dataPreviousTyped.TotalCurrent

	return dataCurrentTyped, nil
}

func (r *dashboardRepository) newDashboardReportProcessor(
	merchantId, period, cacheKeyMask string,
	status interface{},
) (DashboardProcessorRepositoryInterface, error) {
	current := time.Now()
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		return nil, err
	}

	processor := &DashboardReportProcessor{
		Match:       bson.M{"merchant_id": merchantOid, "status": status, "type": "order"},
		Db:          r.db,
		Collection:  CollectionOrderView,
		Cache:       r.cache,
		CacheExpire: time.Duration(0),
	}

	switch period {
	case pkg.DashboardPeriodCurrentDay:
		processor.GroupBy = dashboardReportGroupByHour
		gte := now.BeginningOfDay()
		lte := now.EndOfDay()
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		break
	case pkg.DashboardPeriodPreviousDay, pkg.DashboardPeriodTwoDaysAgo:
		decrement := -1
		if period == pkg.DashboardPeriodTwoDaysAgo {
			decrement = decrement * 2
		}
		previousDay := time.Now().AddDate(0, 0, decrement)
		gte := now.New(previousDay).BeginningOfDay()
		lte := now.New(previousDay).EndOfDay()

		processor.GroupBy = dashboardReportGroupByHour
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		processor.CacheExpire = now.New(current).EndOfDay().Sub(current)
		break
	case pkg.DashboardPeriodCurrentWeek:
		processor.GroupBy = dashboardReportGroupByPeriodInDay
		gte := now.BeginningOfWeek()
		lte := now.EndOfWeek()
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		break
	case pkg.DashboardPeriodPreviousWeek, pkg.DashboardPeriodTwoWeeksAgo:
		decrement := -7
		if period == pkg.DashboardPeriodTwoWeeksAgo {
			decrement = decrement * 2
		}
		previousWeek := time.Now().AddDate(0, 0, decrement)
		gte := now.New(previousWeek).BeginningOfWeek()
		lte := now.New(previousWeek).EndOfWeek()

		processor.GroupBy = dashboardReportGroupByPeriodInDay
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		processor.CacheExpire = now.New(current).EndOfWeek().Sub(current)
		break
	case pkg.DashboardPeriodCurrentMonth:
		gte := now.BeginningOfMonth()
		lte := now.EndOfMonth()

		processor.GroupBy = dashboardReportGroupByDay
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		break
	case pkg.DashboardPeriodPreviousMonth, pkg.DashboardPeriodTwoMonthsAgo:
		decrement := -1
		if period == pkg.DashboardPeriodTwoMonthsAgo {
			decrement = decrement * 2
		}
		previousMonth := now.New(time.Now()).BeginningOfMonth().AddDate(0, decrement, 0)
		gte := now.New(previousMonth).BeginningOfMonth()
		lte := now.New(previousMonth).EndOfMonth()

		processor.GroupBy = dashboardReportGroupByDay
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}

		processor.CacheExpire = now.New(current).EndOfMonth().Sub(current)
		break
	case pkg.DashboardPeriodCurrentQuarter:
		gte := now.BeginningOfQuarter()
		lte := now.EndOfQuarter()
		processor.GroupBy = dashboardReportGroupByWeek
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		break
	case pkg.DashboardPeriodPreviousQuarter, pkg.DashboardPeriodTwoQuarterAgo:
		decrement := -1
		if period == pkg.DashboardPeriodTwoQuarterAgo {
			decrement = decrement * 4
		}
		previousQuarter := now.BeginningOfQuarter().AddDate(0, decrement, 0)
		gte := now.New(previousQuarter).BeginningOfQuarter()
		lte := now.New(previousQuarter).EndOfQuarter()

		processor.GroupBy = dashboardReportGroupByWeek
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}

		processor.CacheExpire = now.New(current).EndOfQuarter().Sub(current)
		break
	case pkg.DashboardPeriodCurrentYear:
		gte := now.BeginningOfYear()
		lte := now.EndOfYear()

		processor.GroupBy = dashboardReportGroupByMonth
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		break
	case pkg.DashboardPeriodPreviousYear, pkg.DashboardPeriodTwoYearsAgo:
		decrement := -1
		if period == pkg.DashboardPeriodTwoYearsAgo {
			decrement = decrement * 2
		}
		previousYear := time.Now().AddDate(decrement, 0, 0)
		gte := now.New(previousYear).BeginningOfYear()
		lte := now.New(previousYear).EndOfYear()

		processor.GroupBy = dashboardReportGroupByMonth
		processor.Match["pm_order_close_date"] = bson.M{"$gte": gte, "$lte": lte}
		processor.CacheExpire = now.New(current).EndOfYear().Sub(current)
		break
	default:
		return nil, fmt.Errorf("incorrect dashboard period")
	}

	if processor.CacheExpire > 0 {
		b, err := json.Marshal(processor.Match)

		if err != nil {
			zap.L().Error(
				"Generate dashboard report cache key failed",
				zap.Error(err),
				zap.Any(pkg.ErrorDatabaseFieldQuery, processor.Match),
			)

			return nil, err
		}

		processor.CacheKey = fmt.Sprintf(cacheKeyMask, md5.Sum(b))
	}

	return processor, nil
}
