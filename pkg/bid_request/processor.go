package bid_request

import (
	"Practice/pkg/storage"
	"github.com/risecodes/openrtb/openrtb2"
	"log"
)

func incrementRequestCount(bidRequest openrtb2.BidRequest) {
	var key string

	for _, imp := range bidRequest.Imp {
		if imp.Banner != nil {
			key = "banner_request_count"
		} else if imp.Video != nil {
			key = "video_request_count"
		}

		if key != "" {
			err := storage.Rdb.Incr(storage.Ctx, key).Err()
			if err != nil {
				log.Printf("Failed to increment %s: %v", key, err)
				continue
			}

			count, err := storage.Rdb.Get(storage.Ctx, key).Int()
			if err != nil {
				log.Printf("Failed to get %s: %v", key, err)
				continue
			}

			log.Printf("Number of %s has been updated: %d", key, count)
		}
	}
}
