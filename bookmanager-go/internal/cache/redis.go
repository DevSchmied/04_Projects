package cache

import (
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the go-redis client.
type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{Client: rdb}
}
