package main

import (
	"log"
	"net/http"

	"Practice/pkg/bid_request"
	"Practice/pkg/rate_limiter"
	"Practice/pkg/storage"
)

func main() {
	log.Println("Initializing...")

	bid_request.InitLogger()
	storage.InitRedis()
	rate_limiter.InitRateLimiter()

	log.Println("Listening on port 8080...")

	go bid_request.ConsumeBidRequests()

	http.HandleFunc("/bid_request", bid_request.HandleBidRequest)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("Server error: %v", err)
	}
}
