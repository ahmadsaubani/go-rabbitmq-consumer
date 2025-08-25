package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis() error {
	fmt.Println("start init redis")
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		fmt.Println("REDIS_HOST environment variable not set")
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		fmt.Println("REDIS_PASSWORD environment variable not set")
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("REDIS_DB environment variable not set")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       redisDB,
	})

	fmt.Println("redis started")

	return nil
}

func SetKey(key string, value interface{}, expiration time.Duration) error {
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func GetKey(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

func DelKey(key string) error {
	return redisClient.Del(ctx, key).Err()
}
