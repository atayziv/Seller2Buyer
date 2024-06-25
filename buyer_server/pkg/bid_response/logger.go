package bid_response

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/risecodes/openrtb/openrtb2"
)

var (
	logFile *os.File
	mu      sync.Mutex
)

func InitLogger() {
	var err error
	logFile, err = os.OpenFile("buyer_responses.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}

	log.Println("Buyer logger initialized.")
}

func LogBidResponse(bidResponse openrtb2.BidResponse) {
	mu.Lock()
	defer mu.Unlock()

	log.Println("Logging the bid response into the log file.")

	data, err := json.Marshal(bidResponse)
	if err != nil {
		log.Printf("Failed to marshal log data: %v", err)
		return
	}

	if _, err := logFile.Write(append(data, '\n')); err != nil {
		log.Printf("Failed to write log data: %v", err)
	}
}
