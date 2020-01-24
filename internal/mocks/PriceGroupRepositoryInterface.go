// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import billingpb "github.com/paysuper/paysuper-proto/go/billingpb"
import context "context"
import mock "github.com/stretchr/testify/mock"

// PriceGroupRepositoryInterface is an autogenerated mock type for the PriceGroupRepositoryInterface type
type PriceGroupRepositoryInterface struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: _a0
func (_m *PriceGroupRepositoryInterface) GetAll(_a0 context.Context) ([]*billingpb.PriceGroup, error) {
	ret := _m.Called(_a0)

	var r0 []*billingpb.PriceGroup
	if rf, ok := ret.Get(0).(func(context.Context) []*billingpb.PriceGroup); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*billingpb.PriceGroup)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: _a0, _a1
func (_m *PriceGroupRepositoryInterface) GetById(_a0 context.Context, _a1 string) (*billingpb.PriceGroup, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.PriceGroup
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.PriceGroup); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PriceGroup)
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

// GetByRegion provides a mock function with given fields: _a0, _a1
func (_m *PriceGroupRepositoryInterface) GetByRegion(_a0 context.Context, _a1 string) (*billingpb.PriceGroup, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *billingpb.PriceGroup
	if rf, ok := ret.Get(0).(func(context.Context, string) *billingpb.PriceGroup); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*billingpb.PriceGroup)
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
func (_m *PriceGroupRepositoryInterface) Insert(_a0 context.Context, _a1 *billingpb.PriceGroup) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PriceGroup) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MultipleInsert provides a mock function with given fields: _a0, _a1
func (_m *PriceGroupRepositoryInterface) MultipleInsert(_a0 context.Context, _a1 []*billingpb.PriceGroup) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*billingpb.PriceGroup) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *PriceGroupRepositoryInterface) Update(_a0 context.Context, _a1 *billingpb.PriceGroup) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *billingpb.PriceGroup) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}