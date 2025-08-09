package redis

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

// Redis tidak ditangani di handler atau main.go.
// Redis ditangani di layer service, agar bersih, reusable, dan scalable.
// Handler hanya terima request (dalam []byte), ubah ke struct, dan kirim ke service.
func InitRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
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
