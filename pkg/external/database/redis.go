package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConnectionOptions struct {
	Prefix   string
	Port     string
	Host     string
	User     string
	Password string
	DB       int
}

func CreateRedisConnection(opt RedisConnectionOptions, logger *zap.Logger) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Username: opt.User,
		Password: opt.Password,
		DB:       opt.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("Redis connection failed", zap.Error(err))
		return
	}

	logger.Info("Redis Connected successfully")

	return
}
