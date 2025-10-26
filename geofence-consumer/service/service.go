package service

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external"
	geofence "github.com/elzestia/fleet/service/geo_fence"
	"github.com/elzestia/fleet/service/healthz"
	"go.uber.org/zap"
)

type Services struct {
	HealthzService  *healthz.HealthzService
	GeoFenceService *geofence.GeoFenceService
}

func CreateServices(cfg *configs.Config, logger *zap.Logger, externalDependencies *external.ExternalDependencies) *Services {
	healthzService := healthz.CreateHalthzService(externalDependencies.PostgreSQLPool, externalDependencies.RedisClient.Client)
	return &Services{
		HealthzService:  healthzService,
		GeoFenceService: geofence.CreateGeoFenceService(logger),
	}
}
