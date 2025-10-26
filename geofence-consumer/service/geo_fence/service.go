package geofence

import "go.uber.org/zap"

type GeoFenceService struct {
	logger *zap.Logger
}

func CreateGeoFenceService(logger *zap.Logger) *GeoFenceService {
	return &GeoFenceService{
		logger: logger,
	}
}
