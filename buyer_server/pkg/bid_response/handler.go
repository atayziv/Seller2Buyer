package bid_response

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/risecodes/openrtb/openrtb2"
)

var (
	bidResponseCh chan openrtb2.BidResponse
)

func init() {
	bidResponseCh = make(chan openrtb2.BidResponse, 100)
}

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

	w.WriteHeader(http.StatusOK)
}

func generatePrice(price float64) float64 {
	return rand.Float64() + price
}

func ProduceBidResponses() {
	for bidResponse := range bidResponseCh {
		go LogBidResponse(bidResponse)
		go ReturnBidResponse(bidResponse)
	}
}
