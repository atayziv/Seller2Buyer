package main

import (
	"Practice/buyer_server/pkg/bid_response"
	"Practice/buyer_server/pkg/server"
)

func main() {
	bid_response.InitLogger()

	go bid_response.ProduceBidResponses()

	s := server.NewServer()
	s.Start()
}
