package main

import (
	"Practice/pkg/bid_request"
	"Practice/pkg/rate_limiter"
	"Practice/pkg/server"
	"Practice/pkg/storage"
)

func main() {
	bid_request.InitLogger()
	storage.InitRedis()
	rate_limiter.InitRateLimiter()

	go bid_request.ConsumeBidRequests()

	s := server.NewServer()
	s.Start()
}
