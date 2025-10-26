package service

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external"
	"github.com/elzestia/fleet/pkg/repository"
	"github.com/elzestia/fleet/service/healthz"
	"github.com/elzestia/fleet/service/user"
	"github.com/elzestia/fleet/service/vehicle"
	"go.uber.org/zap"
)

type Services struct {
	HealthzService *healthz.HealthzService
	UserService    *user.UserService
	VehicleService *vehicle.VehicleService
}

func CreateServices(cfg *configs.Config, logger *zap.Logger, externalDependencies *external.ExternalDependencies) *Services {
	healthzService := healthz.CreateHalthzService(externalDependencies.PostgreSQLPool, externalDependencies.RedisClient.Client)
	repo, err := repository.CreateRepository(externalDependencies.PostgreSQLPool, externalDependencies.RedisClient.Client, logger)
	if err != nil {
		logger.Fatal("failed to create repository", zap.Error(err))
	}

	return &Services{
		HealthzService: healthzService,
		UserService:    user.CreateService(cfg, logger, repo, externalDependencies.RedisClient.Client),
		VehicleService: vehicle.CreateVehicleService(cfg, logger, repo, externalDependencies.RedisClient),
	}
}
