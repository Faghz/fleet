package database

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
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

type RedisClient struct {
	Client *goredislib.Client
	Mutex  *redsync.Redsync
	Prefix string
}

func CreateRedisConnection(opt RedisConnectionOptions, logger *zap.Logger) (client *RedisClient) {
	redisClient := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%s", opt.Host, opt.Port),
		Username: opt.User,
		Password: opt.Password,
		DB:       opt.DB,
	})

	// Test connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("failed to connect to Redis", zap.Error(err))
	}

	pool := goredis.NewPool(redisClient)
	rs := redsync.New(pool)

	client = &RedisClient{
		Client: redisClient,
		Mutex:  rs,
		Prefix: opt.Prefix,
	}

	return
}
