package repository

import "time"

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	IsExist(key string) (bool, error)
}
