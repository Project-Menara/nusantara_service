package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var rdb *redis.Client

func InitRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	log.Println("Redis Connected!")
	return rdb
}

func SetRedis(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return rdb.Set(ctx, key, value, duration).Err()
}

// Get value berdasarkan key
func GetRedis(ctx context.Context, key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Delete key dari Redis
func DeleteRedis(ctx context.Context, key string) error {
	return rdb.Del(ctx, key).Err()
}
