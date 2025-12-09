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
	bookListKey = "book:list"      // Redis key under which book list is stored
	bookListTTL = 60 * time.Second // TTL ensures automatic cache expiration
)

// RedisBookCacher implements BookCacher
type RedisBookCacher struct {
	rdb *RedisClient
}

// The compile-time assertion ensures that RedisBookCacher always satisfies BookCacher.
var _ BookCacher = (*RedisBookCacher)(nil)

// NewRedisBookCache returns a cache instance
func NewRedisBookCacher(rdb *RedisClient) *RedisBookCacher {
	return &RedisBookCacher{rdb: rdb}
}

// GetBookList retrieves the cached book list.
func (c *RedisBookCacher) GetBookList(ctx context.Context) ([]model.Book, error) {
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

// SetBookList stores the book list in Redis with TTL.
func (c *RedisBookCacher) SetBookList(ctx context.Context, books []model.Book) error {
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

// InvalidateBookList removes the cached book list from Redis.
// Called after Add / Update / Delete operations.
func (c *RedisBookCacher) InvalidateBookList(ctx context.Context) error {
	if err := c.rdb.Client.Del(ctx, bookListKey).Err(); err != nil {
		log.Printf("Redis error on Del book:list: %v\n", err)
		return err
	}

	log.Println("Redis Del: book:list deleted")
	return nil
}
