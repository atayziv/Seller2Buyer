package main

import (
	"log"

	"Practice/buyer_server/pkg/bid_response"
	"Practice/buyer_server/pkg/server"
)

func main() {
	log.Println("Initializing Buyer Server...")

	// Initialize logger
	bid_response.InitLogger()
	// Start processing bid responses
	go bid_response.ProduceBidResponses()
	// Initialize and start the server
	s := server.NewServer()
	s.Start()
}
