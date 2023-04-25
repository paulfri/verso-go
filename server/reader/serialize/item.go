package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"gopkg.in/guregu/null.v4"
)

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
			ID:              fmt.Sprintf("%d", item.ID),
			TimestampUsec:   published * 10_000,
			DirectStreamIds: []string{}, // TODO
		}
	})
}

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

func FeedItemsFromReaderIDsRows(rows []query.GetQueueItemsByReaderIDsRow) []FeedItem {
	serializable := lop.Map(rows, func(row query.GetQueueItemsByReaderIDsRow, _ int) SerializableItem {
		return QueueItemByReaderIDsRowToSerializableItem(row)
	})

	return FeedItemsFromSerializable(serializable)
}

func FeedItemsFromUserIDRow(rows []query.GetQueueItemsByUserIDRow) []FeedItem {
	serializable := lop.Map(rows, func(row query.GetQueueItemsByUserIDRow, _ int) SerializableItem {
		return QueueItemByUserIDRowToSerializableItem(row)
	})

	return FeedItemsFromSerializable(serializable)
}

func FeedItemsFromSerializable(items []SerializableItem) []FeedItem {
	return lop.Map(items, func(item SerializableItem, _ int) FeedItem {
		i := item.FeedItem

		published := i.CreatedAt.Unix()
		if i.PublishedAt.Valid {
			published = i.PublishedAt.Time.Unix()
		}

		updated := i.CreatedAt.Unix()
		if i.RemoteUpdatedAt.Valid {
			updated = i.RemoteUpdatedAt.Time.Unix()
		}

		return FeedItem{
			Origin: Origin{
				StreamID: common.ReaderStreamIDFromFeedURL(item.RSSFeedURL),
			},
			ID:     common.LongItemID(i.ReaderID),
			Author: null.String{i.Author},
			Content: FeedItemContent{
				Direction: "ltr",
				Content:   i.Content,
			},
			Summary: FeedItemContent{
				Direction: "ltr",
				Content:   i.Content,
			},
			TimestampUsec: fmt.Sprintf("%d", published*10_000),
			CrawlTimeMsec: fmt.Sprintf("%d", i.CreatedAt.Unix()*1000),
			Published:     published,
			Updated:       updated,
			Title:         i.Title,

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
