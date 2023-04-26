package serialize

import (
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
		published := item.CreatedAt.UnixNano()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.UnixNano()
		}

		return FeedItemRef{
			ID:            common.ShortIDFromReaderID(item.ReaderID),
			TimestampUsec: published,
			DirectStreamIds: []string{
				common.StreamIDReadingList,
				common.ReaderStreamIDFromFeedURL(item.RSSFeedURL),
			},
		}
	})
}
