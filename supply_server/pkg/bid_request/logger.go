package bid_request

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
	logFile, err = os.OpenFile("bid_requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}

	log.Println("Logger initialized.")
}

func logBidRequest(bidRequest openrtb2.BidRequest) {
	mu.Lock()
	defer mu.Unlock()

	log.Println("Logging the bid request into the log file.")

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
