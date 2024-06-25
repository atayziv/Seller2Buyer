package bid_request

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"Practice/supply_server/pkg/storage"

	"github.com/go-redis/redis_rate/v9"
	"github.com/risecodes/openrtb/openrtb2"
)

var (
	bidRequestCh = make(chan openrtb2.BidRequest)
)

func HandleBidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling a new bid request!")

	if !checkRateLimit(w) {
		return
	}

	bidRequest, err := decodeBidRequest(w, r)
	if err != nil {
		return
	}

	log.Printf("Processing bid request for Site ID: %s, Floor price: %.2f", bidRequest.Site.ID, bidRequest.Imp[0].BidFloor)

	bidResponse, err := send2Buyer(bidRequest)
	if err != nil {
		http.Error(w, "Failed to get bid response from buyer", http.StatusInternalServerError)
		return
	}

	bidRequestCh <- *bidRequest

	handleBidResponse(w, bidResponse)
}

func checkRateLimit(w http.ResponseWriter) bool {
	limit := redis_rate.PerMinute(3)
	res, err := storage.RateLimiter.Allow(context.Background(), "global", limit)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return false
	}
	if res.Allowed == 0 {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		log.Println("Too many requests - rate limit exceeded.")
		return false
	}
	return true
}

func decodeBidRequest(w http.ResponseWriter, r *http.Request) (*openrtb2.BidRequest, error) {
	var bidRequest openrtb2.BidRequest
	if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return nil, err
	}
	return &bidRequest, nil
}

func handleBidResponse(w http.ResponseWriter, bidResponse *openrtb2.BidResponse) {
	responseData, err := json.Marshal(bidResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseData); err != nil {
		log.Printf("Error writing bid response: %v", err)
	}
}

func ConsumeBidRequests() {
	for bidRequest := range bidRequestCh {
		go logBidRequest(bidRequest)
		go incrementRequestCount(bidRequest)
	}
}
