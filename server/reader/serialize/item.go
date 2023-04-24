package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"gopkg.in/guregu/null.v4"
)

type FeedItemContent struct {
	Direction string `json:"direction"`
	Content   string `json:"content"`
}

type Origin struct {
	StreamID string `json:"streamId"`
}
type FeedItem struct {
	Origin        Origin          `json:"origin"`
	Updated       int64           `json:"updated"`
	ID            string          `json:"id"`
	Author        null.String     `json:"author"`
	Content       FeedItemContent `json:"content"`
	Summary       FeedItemContent `json:"summary"`
	TimestampUsec string          `json:"timestampUsec"`
	CrawlTimeMsec string          `json:"crawlTimeMsec"`
	Published     int64           `json:"published"`
	Title         string          `json:"title"`

	// TODO:
	Categories  []Category `json:"categories"`
	Alternate   []Category `json:"alternate"`
	Comments    []Category `json:"comments"`
	Annotations []Category `json:"annotations"`
	LikingUsers []Category `json:"likingUsers"`
	Enclosure   []Category `json:"enclosure"`
	MediaGroup  Category   `json:"mediaGroup"`
}

func FeedItemsFromRows(items []query.RSSItem) []FeedItem {
	return lop.Map(items, func(item query.RSSItem, _ int) FeedItem {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		updated := item.CreatedAt.Unix()
		if item.RemoteUpdatedAt.Valid {
			updated = item.RemoteUpdatedAt.Time.Unix()
		}

		return FeedItem{
			Origin: Origin{
				// StreamID: fmt.Sprintf("feed/%d", item.FeedID), // TODO url
				StreamID: "feed/https://www.sounderatheart.com/rss/current.xml",
			},
			ID:     common.LongItemID(item.ReaderID),
			Author: null.String{item.Author},
			Content: FeedItemContent{
				Direction: "ltr",
				Content:   item.Content,
			},
			Summary: FeedItemContent{
				Direction: "ltr",
				Content:   item.Content,
			},
			TimestampUsec: fmt.Sprintf("%d", published*10_000),
			CrawlTimeMsec: fmt.Sprintf("%d", item.CreatedAt.Unix()*1000),
			Published:     published,
			Updated:       updated,
			Title:         item.Title,

			// TODO
			Categories:  []Category{},
			Alternate:   []Category{},
			Comments:    []Category{},
			Annotations: []Category{},
			LikingUsers: []Category{},
			Enclosure:   []Category{},
			MediaGroup:  Category{},
		}
	})
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
