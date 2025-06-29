package redis

import (
	"github.com/redis/go-redis/v9"
)

type (
	RedisRepo struct {
		cache *redis.Client
	}
)

func NewRepo(cache *redis.Client) *RedisRepo {
	return &RedisRepo{
		cache: cache,
	}
}
