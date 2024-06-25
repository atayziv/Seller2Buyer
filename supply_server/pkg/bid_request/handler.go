package bid_request

import (
	"encoding/json"
	"log"
	"net/http"

	"Practice/supply_server/pkg/storage"

	"github.com/go-redis/redis_rate/v9"
	"github.com/risecodes/openrtb/openrtb2"
)

var (
	bidRequestCh chan openrtb2.BidRequest
)

func init() {
	bidRequestCh = make(chan openrtb2.BidRequest, 100)
}

func HandleBidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling a new bid request!")

	limit := redis_rate.PerMinute(3)
	res, err := storage.RateLimiter.Allow(storage.Ctx, "global", limit)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if res.Allowed == 0 {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		log.Println("Too many requests - rate limit exceeded.")
		return
	}

	var bidRequest openrtb2.BidRequest
	if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Processing bid request for Site ID: %s, Floor price: %f", bidRequest.Site.ID, bidRequest.Imp[0].BidFloor)

	bidRequestCh <- bidRequest

	w.WriteHeader(http.StatusOK)
}

func HandleBidResponse(w http.ResponseWriter, r *http.Request) {
	var bidResponse openrtb2.BidResponse
	if err := json.NewDecoder(r.Body).Decode(&bidResponse); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Printf("Received bid response with ID: %s, Buyer offer: %f", bidResponse.ID, bidResponse.SeatBid[0].Bid[0].Price)

	responseData, err := json.Marshal(bidResponse)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set content type and write response data
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
		go send2Buyer(bidRequest)
	}
}
