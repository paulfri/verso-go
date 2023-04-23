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
