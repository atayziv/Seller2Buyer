package bid_request

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/risecodes/openrtb/openrtb2"
)

func send2Buyer(bidRequest *openrtb2.BidRequest) (*openrtb2.BidResponse, error) {
	requestBody, err := json.Marshal(bidRequest)
	if err != nil {
		log.Printf("Failed to marshal bid request: %v", err)
		return nil, err
	}
	resp, err := http.Post("http://localhost:8081/bid_response", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Failed to send bid request to buyer: %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response from buyer server: %v", resp.StatusCode)
	}
	var bidResponse openrtb2.BidResponse
	if err := json.NewDecoder(resp.Body).Decode(&bidResponse); err != nil {
		log.Printf("Failed to decode bid response: %v", err)
		return nil, err
	}

	return &bidResponse, nil
}
