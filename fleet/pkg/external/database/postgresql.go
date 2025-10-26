package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgreSQLConnection struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// CreatePostgreSQLConnection returns a PostgreSQL database connection pool instance using pgx
func CreatePostgreSQLConnection(connConfig PostgreSQLConnection, logger *zap.Logger) *pgxpool.Pool {
	// PostgreSQL DSN format: postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		connConfig.User, connConfig.Password, connConfig.Host, connConfig.Port, connConfig.DbName)

	// Configure connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Fatal("Failed to parse PostgreSQL config", zap.Error(err))
		return nil
	}

	// Set pool configuration
	config.MaxConns = connConfig.MaxConns
	config.MinConns = connConfig.MinConns
	config.MaxConnLifetime = connConfig.MaxConnLifetime
	config.MaxConnIdleTime = connConfig.MaxConnIdleTime

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Fatal("Failed to create PostgreSQL connection pool", zap.Error(err))
		return nil
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		logger.Fatal("Failed to ping PostgreSQL database", zap.Error(err))
		return nil
	}

	logger.Info("PostgreSQL Connected successfully")

	return pool
}
