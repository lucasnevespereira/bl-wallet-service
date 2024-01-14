package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type ICache interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}

type Config struct {
	Host string
	Port int
}

type Cache struct {
	client *redis.Client
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *Cache) Get(key string) ([]byte, error) {
	bytes, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, errors.Errorf("key %s not found", key)
	}
	return bytes, nil
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func NewCacheClient(config Config) ICache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "",
		DB:       0,
	})
	return &Cache{client: rdb}
}
