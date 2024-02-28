package db

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func CreateRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_SEARCH_CACHE"),
		DB:   2,
	})
}
