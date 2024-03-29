// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aholake/order-service/internal/application/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// APIPort is an autogenerated mock type for the APIPort type
type APIPort struct {
	mock.Mock
}

// PlaceOrder provides a mock function with given fields: _a0, _a1
func (_m *APIPort) PlaceOrder(_a0 context.Context, _a1 domain.Order) (domain.Order, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for PlaceOrder")
	}

	var r0 domain.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Order) (domain.Order, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.Order) domain.Order); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(domain.Order)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.Order) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAPIPort creates a new instance of APIPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPIPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *APIPort {
	mock := &APIPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
