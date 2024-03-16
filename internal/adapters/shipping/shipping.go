package shipping

import (
	"context"
	"log"
	"time"

	"github.com/aholake/order-proto/golang/shipping"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ShippingAdapter struct {
	client shipping.ShippingClient
}

func NewShippingAdapter(shippingServiceUrl string) (*ShippingAdapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(circuitBreakerInterceptor()))

	conn, err := grpc.Dial(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	// defer conn.Close()
	client := shipping.NewShippingClient(conn)
	return &ShippingAdapter{
		client: client,
	}, nil
}

func (s ShippingAdapter) Ship(ctx context.Context, orderId int64, address string) error {
	_, err := s.client.Create(ctx, &shipping.CreateShippingRequest{
		Address: address,
	})
	return err
}

func createCb() *gobreaker.CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "shipping-cb",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			log.Printf("total failures: %d, total requests: %d", counts.TotalFailures, counts.Requests)
			return float64(counts.TotalFailures)/float64(counts.Requests) >= 0.6
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("the circuit breaker %s is changed state from %s to %s", name, from, to)
		},
	})
	return cb
}

var cb *gobreaker.CircuitBreaker = createCb()

func circuitBreakerInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		_, cbErr := cb.Execute(func() (interface{}, error) {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
		return cbErr
	}
}
