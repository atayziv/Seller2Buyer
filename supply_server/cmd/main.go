package main

import (
	bidrequest2 "Practice/supply_server/pkg/bid_request"
	"Practice/supply_server/pkg/rate_limiter"
	"Practice/supply_server/pkg/server"
	"Practice/supply_server/pkg/storage"
)

func main() {
	bidrequest2.InitLogger()
	storage.InitRedis()
	rate_limiter.InitRateLimiter()

	go bidrequest2.ConsumeBidRequests()

	s := server.NewServer()
	s.Start()
}
