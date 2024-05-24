package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

func (i *RedisRepository) Set(key string, value interface{}, expiration time.Duration) error {
	return i.rdb.Set(context.Background(), key, value, expiration).Err()
}

func (i *RedisRepository) IsExist(key string) (bool, error) {
	exist, err := i.rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return exist == 1, nil
}
