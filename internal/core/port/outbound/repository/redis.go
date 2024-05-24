package repository

import "time"

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	GetString(key string) (string, error)
	IsExist(key string) (bool, error)
}
