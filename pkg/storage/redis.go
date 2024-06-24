package storage

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

var (
	Rdb         *redis.Client
	RateLimiter *redis_rate.Limiter
	Ctx         = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return
	}

	log.Println("Redis initialized.")
}
