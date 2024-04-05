package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetDataSourceURL() string {
	return getEnvironmentValue("DB_URL")
}

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}

func GetPaymentServiceUrl() string {
	return getEnvironmentValue("PAYMENT_SERVICE_URL")
}

func getEnvironmentValue(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return value
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println(".env file loaded successfully")
}

func GetJWTSecret() string {
	return getEnvironmentValue("JWT_SECRET_KEY")
}

func GetJWTIssuer() string {
	return getEnvironmentValue("JWT_ISSUER")
}

func GetJWTExpiration() int64 {
	expirationStr := getEnvironmentValue("JWT_EXPIRATION")
	expiration, err := strconv.ParseInt(expirationStr, 10, 64)

	if err != nil {
		log.Fatalf("JWT expiration: %s is invalid", expirationStr)
	}

	return expiration
}
