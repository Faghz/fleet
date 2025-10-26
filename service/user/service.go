package user

import (
	"context"

	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/models"
	"github.com/elzestia/fleet/pkg/repository"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

//go:generate mockgen -destination=../../pkg/mocks/user_repository_mock.go -package=mocks github.com/elzestia/fleet/service/user UserRepository

type UserRepository interface {
	GetAuthByUserUserID(ctx context.Context, userID string) (models.Auth, error)
	InsertAuth(ctx context.Context, arg models.InsertAuthParams) error

	GetSessionByEntityId(ctx context.Context, arg models.GetSessionByEntityIdParams) (models.Session, error)
	InsertSession(ctx context.Context, arg models.InsertSessionParams) error
	DeleteSessionByID(ctx context.Context, arg models.DeleteSessionByIDParams) error

	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	InsertUser(ctx context.Context, arg models.InsertUserParams) error

	GetSessionCache(ctx context.Context, userId, sessionId string) (models.Session, error)
	SetSessionCache(ctx context.Context, session models.Session) error
	DeleteSessionCache(ctx context.Context, userId, sessionId string) error

	BeginTx(ctx context.Context) (pgx.Tx, error)
	WithTx(tx pgx.Tx) repository.Querier
	CommitTx(tx pgx.Tx) error
	RollbackTx(tx pgx.Tx) error
}

type UserService struct {
	config     *configs.Config
	logger     *zap.Logger
	repository UserRepository
	redis      *redis.Client
}

func CreateService(
	config *configs.Config,
	logger *zap.Logger,
	repo *repository.Repository,
	redis *redis.Client,
) *UserService {
	return &UserService{
		config:     config,
		logger:     logger,
		redis:      redis,
		repository: repo,
	}
}
