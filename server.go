package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/risecodes/openrtb/openrtb2"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	logFile     *os.File
	mu          sync.Mutex
	rdb         *redis.Client
	rateLimiter *redis_rate.Limiter
	ctx         = context.Background()
)

func init() {
	log.Println("Initializing log file and Redis client...")

	var err error
	logFile, err = os.OpenFile("bid_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Ensure Redis connection is successful
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	rateLimiter = redis_rate.NewLimiter(rdb)

	log.Println("Initialization complete.")
}

func main() {
	log.Println("Listening on port 8080...")

	http.HandleFunc("/bid_request", handleBidRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling a new bid request...")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Rate limiting: limit to 50 requests at a time
	res, err := rateLimiter.Allow(ctx, "global", redis_rate.PerMinute(50)) // limit to 50 requests per minute
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

	go logBidRequest(bidRequest)
	go incrementRequestCount()

	w.WriteHeader(http.StatusOK)
}

func logBidRequest(bidRequest openrtb2.BidRequest) {
	mu.Lock()
	defer mu.Unlock()

	log.Println("Logging the bid request into the log file...")

	logData := struct {
		ID       string `json:"id"`
		SiteName string `json:"site_name"`
		ImpID    string `json:"imp_id"`
		IP       string `json:"ip"`
		UA       string `json:"ua"`
		UserID   string `json:"user_id"`
	}{
		ID:       bidRequest.ID,
		SiteName: bidRequest.Site.Name,
		ImpID:    bidRequest.Imp[0].ID,
		IP:       bidRequest.Device.IP,
		UA:       bidRequest.Device.UA,
		UserID:   bidRequest.User.ID,
	}

	data, err := json.Marshal(logData)
	if err != nil {
		log.Printf("Failed to marshal log data: %v", err)
		return
	}

	if _, err := logFile.Write(append(data, '\n')); err != nil {
		log.Printf("Failed to write log data: %v", err)
	}
}

func incrementRequestCount() {
	err := rdb.Incr(ctx, "request_count").Err()
	if err != nil {
		log.Printf("Failed to increment request count: %v", err)
		return
	}

	count, err := rdb.Get(ctx, "request_count").Int()
	if err != nil {
		log.Printf("Failed to get request count: %v", err)
		return
	}

	log.Printf("Number of bid requests has been updated: %d", count)
}
