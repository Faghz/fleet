package vehicle

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
	"go.uber.org/zap"
)

func (s *VehicleService) checkAndPublishNearestPOI(ctx context.Context, vehicleID string, latitude, longitude float64, timestamp int64) error {
	pointOfInterest, err := s.pointOfInterestService.GetNearestPointOfInterests(ctx, latitude, longitude, 50, 1)
	if err != nil {
		s.logger.Error("[checkAndPublishNearestPOI] Failed to get nearest points of interest", zap.String("vehicle_id", vehicleID), zap.Error(err))
		return err
	}

	if len(pointOfInterest) == 0 {
		s.logger.Debug("[checkAndPublishNearestPOI] No nearby points of interest found",
			zap.String("vehicle_id", vehicleID),
			zap.Float64("latitude", latitude),
			zap.Float64("longitude", longitude))
		return nil
	}

	err = s.publishVehicleReachNearestPOIEvent(ctx, vehicleID, latitude, longitude, timestamp)
	if err != nil {
		s.logger.Error("[checkAndPublishNearestPOI] Failed to publish vehicle reach nearest POI event", zap.String("vehicle_id", vehicleID), zap.Error(err))
		return err
	}

	s.logger.Info("[checkAndPublishNearestPOI] Nearest points of interest",
		zap.String("vehicle_id", vehicleID),
		zap.Any("points_of_interest", pointOfInterest))

	return nil
}

func (s *VehicleService) publishVehicleReachNearestPOIEvent(ctx context.Context, vehicleID string, latitude, longitude float64, timestamp int64) error {
	data := models.ReachedNearestPointOfInterestEvent{
		VehicleID: vehicleID,
		Event:     s.config.RabbitMQ.Publisher.GeoFence.Queue,
		Location: models.ReachedNearestPointOfInterestLocation{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Timestamp: timestamp,
	}

	err := s.rabbitMqClient.Publish(ctx, s.config.RabbitMQ.Publisher.GeoFence.Exchange, s.config.RabbitMQ.Publisher.GeoFence.Queue, data)
	if err != nil {
		s.logger.Error("[publishVehicleReachNearestPOIEvent] Failed to publish vehicle reach nearest POI event", zap.String("vehicle_id", vehicleID), zap.Error(err))
		return err
	}

	return nil
}
