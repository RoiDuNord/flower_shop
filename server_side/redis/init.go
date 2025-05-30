package rds

import (
	"context"
	"fmt"
	"log/slog"
	"server/config"

	"github.com/redis/go-redis/v9"
)

func redisClient() (*redis.Client, error) {
	params, err := config.GetRedisParams()
	if err != nil {
		slog.Error("error getting Redis parameters", "error", err)
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", params.Host, params.Port),
		Password: fmt.Sprintf("%s", params.Password),
	})

	return rdb, nil
}

func InitRedis() (*redis.Client, error) {
	rdb, err := redisClient()
	if err != nil {
		slog.Error("error initializing Redis client", "error", err)
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("error connecting to Redis", "error", err)
		return nil, err
	}

	slog.Info("redis client initialized successfully")
	return rdb, nil
}
