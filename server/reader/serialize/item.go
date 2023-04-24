package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
	"gopkg.in/guregu/null.v4"
)

type FeedItemContent struct {
	Direction string `json:"direction"`
	Content   string `json:"content"`
}

type FeedItem struct {
	// TODO:
	//   categories
	//   alternate
	Origin        EmptyObject     `json:"origin"`
	Updated       int64           `json:"updated"`
	ID            string          `json:"id"`
	Author        null.String     `json:"author"`
	Content       FeedItemContent `json:"content"`
	TimestampUsec int64           `json:"timestampUsec"`
	CrawlTimeMsec int64           `json:"crawlTimeMsec"`
	Published     int64           `json:"published"`
	Title         string          `json:"title"`
}

func FeedItemsFromRows(items []query.GetQueueItemsByUserIDRow) []FeedItem {
	return lop.Map(items, func(item query.GetQueueItemsByUserIDRow, _ int) FeedItem {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		updated := item.CreatedAt.Unix()
		if item.RemoteUpdatedAt.Valid {
			updated = item.RemoteUpdatedAt.Time.Unix()
		}

		return FeedItem{
			Origin: EmptyObject{},
			// TODO: Determine compatibility with Reader
			// https://raw.githubusercontent.com/mihaip/google-reader-api/master/wiki/ItemId.wiki
			ID:     fmt.Sprintf("tag:google.com,2005:reader/item/%016d", item.ID),
			Author: null.String{item.Author},
			Content: FeedItemContent{
				Direction: "ltr",
				Content:   item.Content,
			},
			TimestampUsec: published * 10_000,
			CrawlTimeMsec: item.CreatedAt.Unix() * 1000,
			Published:     published,
			Updated:       updated,
			Title:         item.Title,
		}
	})
}

// * `id`: [ItemId] for the item (signed base 10 version).
//   * `timestampUsec`: time in microseconds since the epoch that the item appeared in the direct stream that it was in.
//   * ``directStreamIds`: array of [StreamId]s representing the direct streams that this item came from.
// lo.Map(items, func(item query.GetQueueItemsByUserIDRow, _ int) int64 {
// 	return item.ID
// })

type FeedItemRef struct {
	ID              string   `json:"id"`
	TimestampUsec   int64    `json:"timestampUsec"`
	DirectStreamIds []string `json:"directStreamIds"`
}

func FeedItemRefsFromRows(items []query.GetQueueItemsByUserIDRow) []FeedItemRef {
	return lop.Map(items, func(item query.GetQueueItemsByUserIDRow, _ int) FeedItemRef {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		return FeedItemRef{
			ID:              fmt.Sprintf("%d", item.ID), // TODO: figure out item IDs. signed base 10?
			TimestampUsec:   published * 10_000,
			DirectStreamIds: []string{"user/0/state/com.google/reading-list"}, // TODO: wtf is this?
		}
	})
}
