package bid_request

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/risecodes/openrtb/openrtb2"
)

func send2Buyer(bidRequest openrtb2.BidRequest) {
	requestBody, err := json.Marshal(bidRequest)
	if err != nil {
		log.Printf("Failed to marshal bid request: %v", err)
		return
	}
	resp, err := http.Post("http://localhost:8081/bid_response", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Failed to send bid request to buyer: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response from buyer server: %v", resp.StatusCode)
	}
}
