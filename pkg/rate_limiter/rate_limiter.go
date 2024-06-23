package rate_limiter

import (
	"Practice/pkg/storage"
	"github.com/go-redis/redis_rate/v9"
	"log"
)

func InitRateLimiter() {
	storage.RateLimiter = redis_rate.NewLimiter(storage.Rdb)
	log.Println("Rate limiter initialized.")
}
