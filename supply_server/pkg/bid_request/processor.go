package bid_request

import (
	"context"

	"Practice/supply_server/pkg/storage"

	"log"

	"github.com/risecodes/openrtb/openrtb2"
)

func incrementRequestCount(bidRequest openrtb2.BidRequest) {
	for _, imp := range bidRequest.Imp {
		var key string

		switch {
		case imp.Banner != nil:
			key = "banner_request_count"
		case imp.Video != nil:
			key = "video_request_count"
		default:
			continue
		}

		err := storage.Rdb.Incr(context.Background(), key).Err()
		if err != nil {
			log.Printf("Failed to increment %s: %v", key, err)
			continue
		}

		count, err := storage.Rdb.Get(context.Background(), key).Int()
		if err != nil {
			log.Printf("Failed to get %s: %v", key, err)
			continue
		}

		log.Printf("Number of %s has been updated: %d", key, count)
	}
}
