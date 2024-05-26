package repository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"synapsis-challenge/shared"
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

func (i *RedisRepository) GetString(key string) (string, error) {
	return i.rdb.Get(context.Background(), key).Result()
}

func (i *RedisRepository) GetInt64(key string) (int64, error) {
	r, e := i.rdb.Get(context.Background(), key).Result()
	if e != nil {
		return 0, e
	}

	return shared.StringToInt64(r, 0), nil
}

func (i *RedisRepository) GetObject(key string, obj interface{}) error {
	r, e := i.rdb.Get(context.Background(), key).Result()
	if e != nil {
		return e
	}

	return json.Unmarshal([]byte(r), obj)
}

func (i *RedisRepository) IsExist(key string) (bool, error) {
	exist, err := i.rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}

	return exist == 1, nil
}
