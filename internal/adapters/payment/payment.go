package payment

import (
	"context"

	"github.com/aholake/order-proto/golang/payment"
	"github.com/aholake/order-service/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	// defer conn.Close()
	client := payment.NewPaymentClient(conn)
	return &Adapter{
		client: client,
	}, nil
}

func (a Adapter) Charge(ctx context.Context, order *domain.Order) error {
	request := payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	}

	_, err := a.client.Create(ctx, &request)

	return err
}
