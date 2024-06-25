package rate_limiter

import (
	"log"

	"Practice/supply_server/pkg/storage"

	"github.com/go-redis/redis_rate/v9"
)

func InitRateLimiter() {
	storage.RateLimiter = redis_rate.NewLimiter(storage.Rdb)
	log.Println("Rate limiter initialized.")
}
