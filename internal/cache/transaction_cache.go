package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type ITransactionCache interface {
	GetTransaction(key string) (string, error)
	SetTransaction(key string, value string) error
	RemoveTransaction(key string) error
}

type TransactionCacheConfig struct {
	Host string
	Port int
}

type TransactionCache struct {
	client *redis.Client
}

func (t *TransactionCache) GetTransaction(key string) (string, error) {
	val, err := t.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", errors.New(fmt.Sprintf("%s does not exist", key))
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (t *TransactionCache) SetTransaction(key, value string) error {
	return t.client.Set(context.Background(), key, value, 0).Err()
}

func (t *TransactionCache) RemoveTransaction(key string) error {
	return t.client.Del(context.Background(), key).Err()
}

func NewTransactionCache(config TransactionCacheConfig) ITransactionCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "",
		DB:       0,
	})
	return &TransactionCache{client: rdb}
}
