package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
)

type SerializableItemRef struct {
	ID          string
	CreatedAt   int64
	PublishedAt int64
}

type FeedItemRef struct {
	ID              string   `json:"id"`
	TimestampUsec   int64    `json:"timestampUsec"`
	DirectStreamIds []string `json:"directStreamIds"`
}

func FeedItemRefsFromRows(items []query.RSSItem) []FeedItemRef {
	return lop.Map(items, func(item query.RSSItem, _ int) FeedItemRef {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		return FeedItemRef{
			ID:              fmt.Sprintf("%d", item.ID),
			TimestampUsec:   published * 10_000,
			DirectStreamIds: []string{}, // TODO
		}
	})
}
