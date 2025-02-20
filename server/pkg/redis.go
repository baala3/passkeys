package pkg

import (
	"github.com/redis/go-redis/v9"
)

func GetRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	return redisClient
}
