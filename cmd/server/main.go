package main

// generate swagger docs
//go:generate go run github.com/swaggo/swag/cmd/swag init -d ../../pkg/transport/http -g ../../../cmd/server/main.go

import (
	"log"

	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external"
	"github.com/elzestia/fleet/pkg/logger"
	api "github.com/elzestia/fleet/pkg/transport"
	inthttp "github.com/elzestia/fleet/pkg/transport/http"
	mqtthandler "github.com/elzestia/fleet/pkg/transport/mqtt/handler"
	"github.com/elzestia/fleet/service"
)

// @title           Transportation API Spec
// @description     Transportation API Specification

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization
func main() {
	// Load configuration.
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	// Initialize logger.
	zapLogger, err := logger.InitLogger(cfg.App.Env, cfg.App.LogLevel)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	// Flushes buffer, if any
	defer zapLogger.Sync()

	// connect infrastructure. Redis, Mongo, etc
	externalDependencies := external.CreateExternalDependencies(cfg, zapLogger)
	// create services
	services := service.CreateServices(cfg, zapLogger, externalDependencies)

	// Create API servers.
	httpServer := api.CreateApiServer(cfg, zapLogger, services)
	// Initialize MQTT handler
	mqttHandler := mqtthandler.CreateMqttConsumer(cfg, zapLogger, externalDependencies.MQTTClient, services.VehicleService)

	inthttp.SetupAndServe(httpServer, services, cfg.Http.Port, cfg.App.ContextTimeout, zapLogger)

	zapLogger.Info("Server started successfully")

	// Wait for termination signal (e.g., SIGINT, SIGTERM) to gracefully shut down the server.
	quit := make(chan struct{})
	<-quit

	zapLogger.Info("Shutting down server...")
	// Shutdown MQTT handler
	mqttHandler.Shutdown()

	// Shutdown server.
	err = api.ShutdownServer(httpServer, zapLogger)
	if err != nil {
		return
	}
}
