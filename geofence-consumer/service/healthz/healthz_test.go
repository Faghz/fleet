package healthz

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPgxPool is a mock for pgxpool.Pool
type MockPgxPool struct {
	mock.Mock
}

func (m *MockPgxPool) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRedisClient is a mock for redis client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

func TestHealthzService_Healthz_Success(t *testing.T) {
	// Create a basic test that checks the service structure
	pgxPool := &pgxpool.Pool{}
	redisClient := &redis.Client{}

	service := CreateHalthzService(pgxPool, redisClient)

	assert.NotNil(t, service)
	assert.NotNil(t, service.pgxPool)
	assert.NotNil(t, service.redisClient)
}

func TestHealthzService_checkPostgreSQL_NilConnection(t *testing.T) {
	service := &HealthzService{
		pgxPool: nil,
	}

	err := service.checkPostgreSQL(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "postgresql connection pool is nil")
}

func TestHealthzService_checkRedis_NilConnection(t *testing.T) {
	service := &HealthzService{
		redisClient: nil,
	}

	err := service.checkRedis(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis client is nil")
}
