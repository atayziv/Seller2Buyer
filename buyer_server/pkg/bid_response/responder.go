package bid_response

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/risecodes/openrtb/openrtb2"
)

func ReturnBidResponse(bidResponse openrtb2.BidResponse) {
	responseData, err := json.Marshal(bidResponse)
	if err != nil {
		log.Printf("Failed to marshal bid response: %v", err)
		return
	}

	// Replace "http://supply_server/bid_response" with the actual supply server endpoint
	resp, err := http.Post("http://localhost:8080/bid_response", "application/json", bytes.NewBuffer(responseData))
	if err != nil {
		log.Printf("Failed to send bid response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response from supply server: %v", resp.StatusCode)
	}
}
