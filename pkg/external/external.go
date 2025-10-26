package external

import (
	"github.com/elzestia/fleet/configs"
	"github.com/elzestia/fleet/pkg/external/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ExternalDependencies struct {
	PostgreSQLPool *pgxpool.Pool
	RedisClient    *redis.Client
}

func CreateExternalDependencies(config *configs.Config, logger *zap.Logger) *ExternalDependencies {
	postgresPool := database.CreatePostgreSQLConnection(database.PostgreSQLConnection{
		Host:            config.Database.PostgreSQL.Host,
		Port:            config.Database.PostgreSQL.Port,
		User:            config.Database.PostgreSQL.User,
		Password:        config.Database.PostgreSQL.Password,
		DbName:          config.Database.PostgreSQL.DbName,
		MaxConns:        config.Database.PostgreSQL.MaxConns,
		MinConns:        config.Database.PostgreSQL.MinConns,
		MaxConnLifetime: config.Database.PostgreSQL.MaxConnLifetime,
		MaxConnIdleTime: config.Database.PostgreSQL.MaxConnIdleTime,
	}, logger)

	redisClient := database.CreateRedisConnection(database.RedisConnectionOptions{
		Host:     config.Database.Redis.Host,
		Port:     config.Database.Redis.Port,
		User:     config.Database.Redis.User,
		Password: config.Database.Redis.Password,
		DB:       config.Database.Redis.DB,
		Prefix:   config.App.Name,
	}, logger)

	return &ExternalDependencies{
		PostgreSQLPool: postgresPool,
		RedisClient:    redisClient,
	}
}
