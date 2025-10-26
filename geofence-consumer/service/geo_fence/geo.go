package geofence

import (
	"context"

	"github.com/elzestia/fleet/pkg/transport/rabbitmq/request"
	"go.uber.org/zap"
)

func (s *GeoFenceService) ConsumeGeoFenceEvents(ctx context.Context, event request.ReachedNearestPointOfInterestEvent) error {
	s.logger.Info("[ConsumeGeoFenceEvents] Consuming GeoFence event", zap.Any("event", event))

	return nil
}
