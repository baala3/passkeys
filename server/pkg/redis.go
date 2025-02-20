package pkg

import (
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client


func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		initRedisClient()
	}
	return redisClient
}
