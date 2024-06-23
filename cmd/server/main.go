package main

import (
	"Practice/pkg/bid_request"
	"Practice/pkg/rate_limiter"
	"Practice/pkg/storage"
	"log"
	"net/http"
)

func main() {
	log.Println("Initializing...")

	bid_request.InitLogger()
	storage.InitRedis()
	rate_limiter.InitRateLimiter()

	log.Println("Listening on port 8080...")

	go bid_request.ConsumeBidRequests()

	http.HandleFunc("/bid_request", bid_request.HandleBidRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
