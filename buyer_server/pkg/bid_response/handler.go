package bid_response

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/risecodes/openrtb/openrtb2"
)

var (
	bidResponseCh = make(chan openrtb2.BidResponse)
)

func HandleBidRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Buyer server handling a new bid request!")

	var bidRequest openrtb2.BidRequest
	if err := json.NewDecoder(r.Body).Decode(&bidRequest); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	price := generatePrice(bidRequest.Imp[0].BidFloor)

	bidResponse := openrtb2.BidResponse{
		ID: bidRequest.ID,
		SeatBid: []openrtb2.SeatBid{
			{
				Bid: []openrtb2.Bid{
					{
						ID:    "1",
						ImpID: bidRequest.Imp[0].ID,
						Price: price,
					},
				},
			},
		},
	}

	bidResponseCh <- bidResponse

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

func generatePrice(price float64) float64 {
	return rand.Float64() + price
}

func ProduceBidResponses() {
	for bidResponse := range bidResponseCh {
		go LogBidResponse(&bidResponse)
	}
}
