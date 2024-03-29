// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aholake/order-service/internal/application/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// PaymentClientPort is an autogenerated mock type for the PaymentClientPort type
type PaymentClientPort struct {
	mock.Mock
}

// Charge provides a mock function with given fields: ctx, order
func (_m *PaymentClientPort) Charge(ctx context.Context, order *domain.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Charge")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPaymentClientPort creates a new instance of PaymentClientPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentClientPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentClientPort {
	mock := &PaymentClientPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
