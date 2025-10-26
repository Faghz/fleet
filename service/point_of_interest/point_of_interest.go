package pointofinterest

import (
	"context"

	"github.com/elzestia/fleet/pkg/models"
	"go.uber.org/zap"
)

func (s *PointOfInterestService) SeedRedisPointsOfInterest() error {
	// Implementation to seed Redis with initial points of interest data.
	ctx := context.Background()
	pointOfInterests, err := s.pointOfInterestRepository.GetPointOfInterests(ctx)
	if err != nil {
		s.logger.Error("[SeedRedisPointsOfInterest] failed to get points of interest from repository", zap.Error(err))
		return err
	}

	err = s.pointOfInterestRepository.SetGeoSpatialPointOfInterests(ctx, pointOfInterests)
	if err != nil {
		s.logger.Error("[SeedRedisPointsOfInterest] failed to seed Redis with points of interest", zap.Error(err))
		return err
	}

	return nil
}

func (s *PointOfInterestService) GetNearestPointOfInterests(ctx context.Context, latitude, longitude float64, radiusMeters float64, limit int) ([]models.PointOfInterest, error) {
	points, err := s.pointOfInterestRepository.GetNearestPointOfInterests(ctx, latitude, longitude, radiusMeters, limit)
	if err != nil {
		s.logger.Error("[GetNearestPointOfInterests] failed to get nearest points of interest", zap.Error(err))
		return nil, err
	}

	var response []models.PointOfInterest
	for _, point := range points {
		response = append(response, models.PointOfInterest{
			ID:          point.ID,
			Name:        point.Name,
			Description: point.Description,
			Latitude:    point.Latitude,
			Longitude:   point.Longitude,
			CreatedBy:   point.CreatedBy,
			CreatedAt:   point.CreatedAt,
			UpdatedBy:   point.UpdatedBy,
			UpdatedAt:   point.UpdatedAt,
			DeletedBy:   point.DeletedBy,
			DeletedAt:   point.DeletedAt,
		})
	}

	return response, nil
}
