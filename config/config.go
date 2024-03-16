package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	postStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(postStr)
	if err != nil {
		log.Fatalf("Unable to convert port value: %v", err)
	}
	return port
}

func GetPaymentServiceURL() string {
	return getEnvironmentValue("PAYMENT_SERVICE_URL")
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("The environment key %v is not found", key)
	}
	return os.Getenv(key)
}
