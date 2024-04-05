package main

import (
	"github.com/diofanto33/cocosette-auth/config"
	"github.com/diofanto33/cocosette-auth/internal/adapters/db"
	"github.com/diofanto33/cocosette-auth/internal/adapters/grpc"
	"github.com/diofanto33/cocosette-auth/internal/adapters/jwt"
	"github.com/diofanto33/cocosette-auth/internal/application/core/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	/* Load environment variables */
	config.LoadEnv()

	/* Create database and JWT adapters */

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	jwtAdapter, err := jwt.NewAdapter(config.GetJWTSecret(), config.GetJWTIssuer(), config.GetJWTExpiration())
	if err != nil {
		log.Fatalf("Failed to create JWT adapter. Error: %v", err)
	}

	/* Dependency Injection */
	application := api.NewApplication(dbAdapter, jwtAdapter)

	/* Run gRPC server */
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
