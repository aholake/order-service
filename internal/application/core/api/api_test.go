package api

import (
	"context"
	"testing"

	"github.com/aholake/order-service/internal/application/core/domain"
	"github.com/aholake/order-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Should_Place_Order_Successfully(t *testing.T) {
	payment := mocks.NewPaymentClientPort(t)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)

	db := mocks.NewDBPort(t)
	db.On("Save", mock.Anything).Return(nil)

	shipping := mocks.NewShippingClientPort(t)
	shipping.On("Ship", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	app := NewApplication(db, payment, shipping)
	order, err := app.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 1,
		OrderItems: []domain.OrderItem{},
	})
	assert.NotNil(t, order)
	assert.Nil(t, err)
}
