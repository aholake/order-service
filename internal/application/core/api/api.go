package api

import (
	"context"
	"fmt"

	"github.com/aholake/order-service/internal/adapters/payment"
	"github.com/aholake/order-service/internal/adapters/shipping"
	"github.com/aholake/order-service/internal/application/core/domain"
	"github.com/aholake/order-service/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db       ports.DBPort
	payment  payment.Adapter
	shipping shipping.ShippingAdapter
}

func NewApplication(db ports.DBPort, paymentAdapter payment.Adapter, shipping shipping.ShippingAdapter) *Application {
	return &Application{
		db:       db,
		payment:  paymentAdapter,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	error := a.db.Save(&order)
	if error != nil {
		return domain.Order{}, error
	}
	err := a.payment.Charge(ctx, &order)
	if err != nil {
		s, _ := status.FromError(err)
		details := errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: s.Message(),
		}
		badReq := errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{&details},
		}
		status, _ := status.New(codes.InvalidArgument, "order creation failed").WithDetails(&badReq)
		return domain.Order{}, status.Err()
	}

	err = a.shipping.Ship(ctx, order.ID, fmt.Sprintf("Address for order %d", order.ID))
	if err != nil {
		status := status.New(codes.InvalidArgument, "shipping creation failed")
		return domain.Order{}, status.Err()
	}

	return order, nil
}
