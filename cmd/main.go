package main

import (
	"log"

	"github.com/aholake/order-service/config"
	"github.com/aholake/order-service/internal/adapters/db"
	"github.com/aholake/order-service/internal/adapters/grpc"
	"github.com/aholake/order-service/internal/adapters/payment"
	"github.com/aholake/order-service/internal/adapters/shipping"
	"github.com/aholake/order-service/internal/application/core/api"
)

func main() {

	dbAdapter, error := db.NewAdapter(config.GetDataSourceURL())
	if error != nil {
		log.Fatalf("Failed to connect to DB, error: %v", error)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("unable to connect to payment service, %v", err)
	}
	shippingAdapter, err := shipping.NewShippingAdapter(config.GetShippingServiceURL())
	if err != nil {
		log.Fatalf("unable to connect to shipping service, %v", err)
	}

	application := api.NewApplication(dbAdapter, *paymentAdapter, *shippingAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
