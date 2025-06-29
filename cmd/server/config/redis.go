package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SetupCache(cfg *Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS_HOSTNAME, cfg.REDIS_PORT),
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
