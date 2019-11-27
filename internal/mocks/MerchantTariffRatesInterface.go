// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billing "github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
import context "context"
import grpc "github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
import mock "github.com/stretchr/testify/mock"

// MerchantTariffRatesInterface is an autogenerated mock type for the MerchantTariffRatesInterface type
type MerchantTariffRatesInterface struct {
	mock.Mock
}

// GetBy provides a mock function with given fields: ctx, in
func (_m *MerchantTariffRatesInterface) GetBy(ctx context.Context, in *grpc.GetMerchantTariffRatesRequest) (*grpc.GetMerchantTariffRatesResponseItems, error) {
	ret := _m.Called(ctx, in)

	var r0 *grpc.GetMerchantTariffRatesResponseItems
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) *grpc.GetMerchantTariffRatesResponseItems); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.GetMerchantTariffRatesResponseItems)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCacheKeyForGetBy provides a mock function with given fields: _a0
func (_m *MerchantTariffRatesInterface) GetCacheKeyForGetBy(_a0 *grpc.GetMerchantTariffRatesRequest) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(*grpc.GetMerchantTariffRatesRequest) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*grpc.GetMerchantTariffRatesRequest) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaymentTariffsBy provides a mock function with given fields: _a0, _a1
func (_m *MerchantTariffRatesInterface) GetPaymentTariffsBy(_a0 context.Context, _a1 *grpc.GetMerchantTariffRatesRequest) ([]*billing.MerchantTariffRatesPayment, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*billing.MerchantTariffRatesPayment
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) []*billing.MerchantTariffRatesPayment); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billing.MerchantTariffRatesPayment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTariffsSettings provides a mock function with given fields: ctx, in
func (_m *MerchantTariffRatesInterface) GetTariffsSettings(ctx context.Context, in *grpc.GetMerchantTariffRatesRequest) (*billing.MerchantTariffRatesSettings, error) {
	ret := _m.Called(ctx, in)

	var r0 *billing.MerchantTariffRatesSettings
	if rf, ok := ret.Get(0).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) *billing.MerchantTariffRatesSettings); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billing.MerchantTariffRatesSettings)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *grpc.GetMerchantTariffRatesRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
