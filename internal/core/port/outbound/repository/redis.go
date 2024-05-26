package repository

import "time"

type RedisRepository interface {
	SetObj(key string, value interface{}, expiration time.Duration) error
	SetString(key string, value string, expiration time.Duration) error
	SetInt64(key string, value int64, expiration time.Duration) error
	GetString(key string) (string, error)
	GetInt64(key string) (int64, error)
	GetObject(key string, obj interface{}) error
	IsExist(key string) (bool, error)
}
