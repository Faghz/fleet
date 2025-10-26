package pointofinterest

import (
	"context"

	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/rabbitmq"
	"github.com/elzestia/fleet/pkg/models"
	"go.uber.org/zap"
)

type PointOfInterestRepository interface {
	SetGeoSpatialPointOfInterests(ctx context.Context, points []models.PointOfInterest) error
	GetPointOfInterests(ctx context.Context) ([]models.PointOfInterest, error)
	GetNearestPointOfInterests(ctx context.Context, latitude, longitude float64, radius float64, count int) ([]models.PointOfInterest, error)
}

type PointOfInterestService struct {
	config                    *configs.Config
	rabbitMQClient            *rabbitmq.RabbitMQClient
	pointOfInterestRepository PointOfInterestRepository
	logger                    *zap.Logger
}

func CreatePointOfInterestService(config *configs.Config, logger *zap.Logger, pointOfInterestRepo PointOfInterestRepository, rabbitMQClient *rabbitmq.RabbitMQClient) *PointOfInterestService {
	pointOfInterestService := &PointOfInterestService{
		config:                    config,
		rabbitMQClient:            rabbitMQClient,
		pointOfInterestRepository: pointOfInterestRepo,
		logger:                    logger,
	}

	go pointOfInterestService.SeedRedisPointsOfInterest()

	return pointOfInterestService
}
