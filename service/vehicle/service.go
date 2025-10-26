package vehicle

import (
	"github.com/elzestia/fleet/configs"
	"go.uber.org/zap"
)

type VehicleService struct {
	config *configs.Config
	logger *zap.Logger
}

func CreateVehicleService(config *configs.Config, logger *zap.Logger) *VehicleService {
	return &VehicleService{
		config: config,
		logger: logger,
	}
}
