package services

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

type RedisClientImpl struct {
	client *redis.Client
}

func NewRedisClient() *RedisClientImpl {
	// Create a new Redis client instance here
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return &RedisClientImpl{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf(os.Getenv("REDIS_ADDRESS") + ":" + os.Getenv("REDIS_PORT")),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		}),
	}
}

func (r *RedisClientImpl) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *RedisClientImpl) Get(key string) (interface{}, error) {
	return r.client.Get(key).Result()
}

func (r *RedisClientImpl) Del(key string) error {
	return r.client.Del(key).Err()
}
