package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})
}
