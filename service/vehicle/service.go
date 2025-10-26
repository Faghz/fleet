package vehicle

import (
	"context"

	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/database"
	"github.com/elzestia/fleet/pkg/external/rabbitmq"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/repository"
	pointofinterest "github.com/elzestia/fleet/service/point_of_interest"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type VehicleRepository interface {
	GetVehicleByVehicleID(ctx context.Context, vehicleID string) (models.Vehicle, error)
	InsertVehicleLocation(ctx context.Context, arg models.InsertVehicleLocationParams) (models.VehicleLocation, error)
	GetVehicleLatestLocationByVehicleID(ctx context.Context, vehicleID string) (models.GetVehicleLatestLocationByVehicleIDRow, error)
	GetVehicleLocationHistoryByVehicleIDAndTimeRange(ctx context.Context, arg models.GetVehicleLocationHistoryByVehicleIDParams) ([]models.GetVehicleLocationHistoryByVehicleIDRow, error)

	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) repository.Querier
	CommitTx(tx pgx.Tx) error
	RollbackTx(tx pgx.Tx) error
}

type VehicleService struct {
	config                 *configs.Config
	logger                 *zap.Logger
	repo                   VehicleRepository
	redis                  *database.RedisClient
	rabbitMqClient         *rabbitmq.RabbitMQClient
	pointOfInterestService *pointofinterest.PointOfInterestService
}

func CreateVehicleService(
	config *configs.Config,
	logger *zap.Logger,
	repo VehicleRepository,
	redis *database.RedisClient,
	pointOfInterestService *pointofinterest.PointOfInterestService,
	rabbitMqClient *rabbitmq.RabbitMQClient,
) *VehicleService {
	return &VehicleService{
		config:                 config,
		logger:                 logger,
		repo:                   repo,
		redis:                  redis,
		pointOfInterestService: pointOfInterestService,
		rabbitMqClient:         rabbitMqClient,
	}
}
