package repository

import "time"

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	GetString(key string) (string, error)
	GetInt64(key string) (int64, error)
	GetObject(key string, obj interface{}) error
	IsExist(key string) (bool, error)
}
