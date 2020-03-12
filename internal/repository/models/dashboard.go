package models

import (
	"github.com/paysuper/paysuper-proto/go/billingpb"
	tools "github.com/paysuper/paysuper-tools/number"
)

type dashboardAmountItemWithChartMapper struct{}

func NewDashboardAmountItemWithChartMapper() Mapper {
	return &dashboardAmountItemWithChartMapper{}
}

type MgoGrossRevenueAndVatReports struct {
	GrossRevenue *MgoDashboardAmountItemWithChart
	Vat          *MgoDashboardAmountItemWithChart
}

type MgoDashboardAmountItemWithChart struct {
	Amount   float64                              `bson:"amount"`
	Currency string                               `bson:"currency"`
	Chart    []*billingpb.DashboardChartItemFloat `bson:"chart"`
}

func (m *dashboardAmountItemWithChartMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	panic("Implement me!")
	return nil, nil
}

func (m *dashboardAmountItemWithChartMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	in := obj.(*MgoDashboardAmountItemWithChart)

	out := &billingpb.DashboardAmountItemWithChart{
		AmountCurrent: tools.FormatAmount(in.Amount),
		Currency:      in.Currency,
	}

	for _, v := range in.Chart {
		item := &billingpb.DashboardChartItemFloat{
			Label: v.Label,
			Value: tools.FormatAmount(v.Value),
		}
		out.Chart = append(out.Chart, item)
	}

	return out, nil
}

type dashboardRevenueDynamicReportItemMapper struct{}

func NewDashboardRevenueDynamicReportItemMapper() Mapper {
	return &dashboardRevenueDynamicReportItemMapper{}
}

type MgoDashboardRevenueDynamicReport struct {
	Currency string
	Items    []*MgoDashboardRevenueDynamicReportItem
}

type MgoDashboardRevenueDynamicReportItem struct {
	Label    int64   `bson:"label"`
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
	Count    int64   `bson:"count"`
}

func (m *dashboardRevenueDynamicReportItemMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	panic("Implement me!")
	return nil, nil
}

func (m *dashboardRevenueDynamicReportItemMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	in := obj.(*MgoDashboardRevenueDynamicReportItem)

	out := &billingpb.DashboardRevenueDynamicReportItem{
		Amount:   tools.FormatAmount(in.Amount),
		Label:    in.Label,
		Currency: in.Currency,
		Count:    in.Count,
	}

	return out, nil
}

type dashboardRevenueByCountryReportMapper struct{}

func NewDashboardRevenueByCountryReportMapper() Mapper {
	return &dashboardRevenueByCountryReportMapper{}
}

type MgoDashboardRevenueByCountryReport struct {
	Currency      string                                         `bson:"currency"`
	Top           []*MgoDashboardRevenueByCountryReportTop       `bson:"top"`
	TotalCurrent  float64                                        `bson:"total"`
	TotalPrevious float64                                        `bson:"total_previous"`
	Chart         []*MgoDashboardRevenueByCountryReportChartItem `bson:"chart"`
}

type MgoDashboardRevenueByCountryReportTop struct {
	Country string  `bson:"_id"`
	Amount  float64 `bson:"amount"`
}

type MgoDashboardRevenueByCountryReportChartItem struct {
	Label  int64   `bson:"label"`
	Amount float64 `bson:"amount"`
}

func (m *dashboardRevenueByCountryReportMapper) MapObjectToMgo(obj interface{}) (interface{}, error) {
	panic("Implement me!")
	return nil, nil
}

func (m *dashboardRevenueByCountryReportMapper) MapMgoToObject(obj interface{}) (interface{}, error) {
	in := obj.(*MgoDashboardRevenueByCountryReport)

	out := &billingpb.DashboardRevenueByCountryReport{
		TotalCurrent:  tools.FormatAmount(in.TotalCurrent),
		TotalPrevious: tools.FormatAmount(in.TotalPrevious),
		Currency:      in.Currency,
	}

	for _, item := range in.Top {
		out.Top = append(out.Top, &billingpb.DashboardRevenueByCountryReportTop{
			Amount:  tools.FormatAmount(item.Amount),
			Country: item.Country,
		})
	}

	for _, item := range in.Chart {
		out.Chart = append(out.Chart, &billingpb.DashboardRevenueByCountryReportChartItem{
			Amount: tools.FormatAmount(item.Amount),
			Label:  item.Label,
		})
	}

	return out, nil
}
