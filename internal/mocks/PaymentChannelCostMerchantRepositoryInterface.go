// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"
import pkg "github.com/paysuper/paysuper-billing-server/internal/pkg"

// PaymentChannelCostMerchantRepositoryInterface is an autogenerated mock type for the PaymentChannelCostMerchantRepositoryInterface type
type PaymentChannelCostMerchantRepositoryInterface struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, obj
func (_m *PaymentChannelCostMerchantRepositoryInterface) Delete(ctx context.Context, obj *billingpb.PaymentChannelCostMerchant) error {
	ret := _m.Called(ctx, obj)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PaymentChannelCostMerchant) error); ok {
		r0 = rf(ctx, obj)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *PaymentChannelCostMerchantRepositoryInterface) Find(_a0 context.Context, _a1 string, _a2 string, _a3 string, _a4 string, _a5 string, _a6 string) ([]*pkg.PaymentChannelCostMerchantSet, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 []*pkg.PaymentChannelCostMerchantSet
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, string) []*pkg.PaymentChannelCostMerchantSet); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*pkg.PaymentChannelCostMerchantSet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllForMerchant provides a mock function with given fields: _a0, _a1
func (_m *PaymentChannelCostMerchantRepositoryInterface) GetAllForMerchant(_a0 context.Context, _a1 string) ([]*billingpb.PaymentChannelCostMerchant, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*billingpb.PaymentChannelCostMerchant
	if rf, ok := ret.Get(0).(func(context.Context, string) []*billingpb.PaymentChannelCostMerchant); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.PaymentChannelCostMerchant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: _a0, _a1
func (_m *PaymentChannelCostMerchantRepositoryInterface) GetById(_a0 context.Context, _a1 string) (*billingpb.PaymentChannelCostMerchant, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.PaymentChannelCostMerchant
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.PaymentChannelCostMerchant); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PaymentChannelCostMerchant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *PaymentChannelCostMerchantRepositoryInterface) Insert(_a0 context.Context, _a1 *billingpb.PaymentChannelCostMerchant) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PaymentChannelCostMerchant) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultipleInsert provides a mock function with given fields: _a0, _a1
func (_m *PaymentChannelCostMerchantRepositoryInterface) MultipleInsert(_a0 context.Context, _a1 []*billingpb.PaymentChannelCostMerchant) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*billingpb.PaymentChannelCostMerchant) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *PaymentChannelCostMerchantRepositoryInterface) Update(_a0 context.Context, _a1 *billingpb.PaymentChannelCostMerchant) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PaymentChannelCostMerchant) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}