package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/server/reader/common"
)

type FeedItemRef struct {
	ID              string   `json:"id"`
	TimestampUsec   int64    `json:"timestampUsec"`
	DirectStreamIds []string `json:"directStreamIds"`
}

func FeedItemRefsFromRows(items []SerializableItem) []FeedItemRef {
	return lop.Map(items, func(item SerializableItem, _ int) FeedItemRef {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		return FeedItemRef{
			ID:            fmt.Sprintf("%d", item.ID),
			TimestampUsec: published * 10_000,
			DirectStreamIds: []string{
				common.StreamIDReadingList,
				common.ReaderStreamIDFromFeedURL(item.RSSFeedURL),
			},
		}
	})
}
