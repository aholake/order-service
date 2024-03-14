package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/aholake/order-proto/golang/order"
	"github.com/aholake/order-service/config"
	"github.com/aholake/order-service/internal/application/core/domain"
	"github.com/aholake/order-service/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("Unable to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc on port %d, error: %v", a.port, err)
	}
}

func (a Adapter) Create(context context.Context, orderRequest *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	orderItems := []domain.OrderItem{}
	for _, oi := range orderRequest.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: oi.ProductCode,
			UnitPrice:   oi.UnitPrice,
			Quantity:    oi.Quantity,
		})
	}
	newOrder := domain.NewOrder(orderRequest.UserId, orderItems)
	res, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: res.ID,
	}, nil
}