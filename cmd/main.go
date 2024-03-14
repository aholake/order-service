package main

import (
	"log"

	"github.com/aholake/order-service/config"
	"github.com/aholake/order-service/internal/adapters/db"
	"github.com/aholake/order-service/internal/adapters/grpc"
	"github.com/aholake/order-service/internal/application/core/api"
)

func main() {
	dbAdapter, error := db.NewAdapter(config.GetDataSourceURL())
	if error != nil {
		log.Fatalf("Failed to connect to DB, error: %v", error)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
