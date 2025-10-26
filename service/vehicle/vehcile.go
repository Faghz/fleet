package vehicle

import (
	"github.com/elzestia/fleet/pkg/transport/mqtt/request"
	"go.uber.org/zap"
)

// processVehicleLocation handles the business logic for vehicle location updates
func (h *VehicleService) ProcessVehicleLocationSync(data *request.VehicleLocationRequest) {

	h.logger.Info("Processing vehicle location",
		zap.String("vehicle_id", data.VehicleID),
		zap.Float64("latitude", data.Latitude),
		zap.Float64("longitude", data.Longitude),
		zap.Int64("timestamp", data.Timestamp))
}
