package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
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
	Categories    []string        `json:"categories"`

	// TODO:
	Alternate   []EmptyObject `json:"alternate"`
	Comments    []EmptyObject `json:"comments"`
	Annotations []EmptyObject `json:"annotations"`
	LikingUsers []EmptyObject `json:"likingUsers"`
	Enclosure   []EmptyObject `json:"enclosure"`
	MediaGroup  EmptyObject   `json:"mediaGroup"`
}

func FeedItemsFromSerializable(items []SerializableItem) []FeedItem {
	return lop.Map(items, func(item SerializableItem, _ int) FeedItem {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		updated := item.CreatedAt.Unix()
		if item.RemoteUpdatedAt.Valid {
			updated = item.RemoteUpdatedAt.Time.Unix()
		}

		var categories []string
		if !item.Unread {
			categories = append(categories, common.StreamIDRead)
		}
		if item.Starred {
			categories = append(categories, common.StreamIDStarred)
		}
		if item.UserID > 0 {
			categories = append(categories, common.StreamIDReadingList)
		}

		return FeedItem{
			Origin: Origin{
				StreamID: common.ReaderStreamIDFromFeedURL(item.RSSFeedURL),
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
			Categories:    categories,

			// TODO
			Alternate:   []EmptyObject{},
			Comments:    []EmptyObject{},
			Annotations: []EmptyObject{},
			LikingUsers: []EmptyObject{},
			Enclosure:   []EmptyObject{},
			MediaGroup:  EmptyObject{},
		}
	})
}
