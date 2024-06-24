package bid_request

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis_rate/v9"
	"github.com/risecodes/openrtb/openrtb2"

	"Practice/pkg/storage"
)

var (
	bidRequestCh chan openrtb2.BidRequest
)

func init() {
	bidRequestCh = make(chan openrtb2.BidRequest, 100)
}

func HandleBidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling a new bid request!")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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

	log.Printf("Processing bid request for Site ID: %s, Device IP: %s", bidRequest.Site.ID, bidRequest.Device.IP)

	bidRequestCh <- bidRequest

	w.WriteHeader(http.StatusOK)
}

func ConsumeBidRequests() {
	for bidRequest := range bidRequestCh {
		go logBidRequest(bidRequest)
		go incrementRequestCount(bidRequest)
	}
}
