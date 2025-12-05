package cache

import "time"

const (
	bookListKey = "book:list"
	bookListTTL = 60 * time.Second
)

// implement interface BookCache
type RedisBookCache struct {
	rdb *RedisClient
}

func NewRedisBookCache(rdb *RedisClient) *RedisBookCache {
	return &RedisBookCache{rdb: rdb}
}
