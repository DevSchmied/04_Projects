package cache

import (
	"bookmanager-go/internal/model"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func (c *RedisBookCache) GetBookList(ctx context.Context) ([]model.Book, error) {
	bookList, err := c.rdb.Client.Get(ctx, bookListKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss (redis.Nil) — not an error, return no data without error.
			log.Println("Redis miss: book:list not found")
			return nil, nil
		}
		log.Printf("Redis error on GET book:list: %v\n", err)
		return nil, err
	}

	log.Println("Redis Hit: book:list loaded")

	var books []model.Book
	if err := json.Unmarshal([]byte(bookList), &books); err != nil {
		log.Printf("JSON error: failed to unmarshal book list: %v\n", err)
		return nil, err
	}
	return books, nil
}

func (c *RedisBookCache) SetBookList(ctx context.Context, books []model.Book) error {
	jsonBytes, err := json.Marshal(books)
	if err != nil {
		log.Printf("JSON error: failed to marshal book list: %v\n", err)
		return err
	}

	err = c.rdb.Client.Set(ctx, bookListKey, jsonBytes, bookListTTL).Err()
	if err != nil {
		log.Printf("Redis error on SET book:list: %v\n", err)
		return err
	}

	log.Println("Redis Set: book:list updated")
	return nil
}
