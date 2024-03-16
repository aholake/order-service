package payment

import (
	"context"
	"time"

	"github.com/aholake/order-proto/golang/payment"
	"github.com/aholake/order-service/internal/application/core/domain"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Millisecond*500)))))

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
