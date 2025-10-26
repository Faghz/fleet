package healthz

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthzService struct {
	pgxPool     *pgxpool.Pool
	redisClient *redis.Client
}

func CreateHalthzService(pgxPool *pgxpool.Pool, redisClient *redis.Client) *HealthzService {
	return &HealthzService{
		pgxPool:     pgxPool,
		redisClient: redisClient,
	}
}

func (s *HealthzService) Healthz(ctx context.Context) (err error) {
	// Check PostgreSQL connection
	if err = s.checkPostgreSQL(ctx); err != nil {
		return fmt.Errorf("postgresql health check failed: %w", err)
	}

	// Check Redis connection
	if err = s.checkRedis(ctx); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

func (s *HealthzService) checkPostgreSQL(ctx context.Context) error {
	if s.pgxPool == nil {
		return fmt.Errorf("postgresql connection pool is nil")
	}

	// Ping the database with context
	if err := s.pgxPool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping postgresql: %w", err)
	}

	return nil
}

func (s *HealthzService) checkRedis(ctx context.Context) error {
	if s.redisClient == nil {
		return fmt.Errorf("redis client is nil")
	}

	// Ping Redis with context
	if err := s.redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	return nil
}
