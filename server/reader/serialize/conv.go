package serialize

import (
	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
)

// TODO: Go generics don't support shared fields so we need all this conversion
// boilerplate. Maybe in 1.22?
// https://github.com/golang/go/issues/48522
//
// type feedItemRows interface {
// 	query.GetQueueItemsByUserIDRow | query.GetQueueItemsByReaderIDsRow
// }
//
// func FeedItemsFromRows[feedItemRows](items []T) []FeedItem {
// 	return lop.Map(items, func(item T, _ int) FeedItem {
// 		published := item.CreatedAt.Unix()
// 		if item.PublishedAt.Valid {
// 			published = item.PublishedAt.Time.Unix()
// 		}
//
// 		...
//
// 		return FeedItem{...}
// 	})
// }
//
//	FeedItemsFromRows[query.GetQueueItemsByUserIDRow](rows)

type SerializableItem struct {
	RSSItem    *query.RSSItem
	RSSFeedURL string
}

func SerializableItemsFromQueueItemByReaderIDsRows(rows []query.GetItemsWithURLByReaderIDsRow) []SerializableItem {
	return lop.Map(rows, func(row query.GetItemsWithURLByReaderIDsRow, _ int) SerializableItem {
		return SerializableItemFromQueueItemByReaderIDsRow(row)
	})
}

func SerializableItemFromQueueItemByReaderIDsRow(item query.GetItemsWithURLByReaderIDsRow) SerializableItem {
	return SerializableItem{
		RSSItem: &query.RSSItem{
			ID:              item.ReaderID,
			UUID:            item.UUID,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
			FeedID:          item.FeedID,
			Title:           item.Title,
			RSSGuid:         item.RSSGuid,
			Link:            item.Link,
			Author:          item.Author,
			AuthorEmail:     item.AuthorEmail,
			Content:         item.Content,
			Summary:         item.Summary,
			PublishedAt:     item.PublishedAt,
			RemoteUpdatedAt: item.RemoteUpdatedAt,
			ReaderID:        item.ReaderID,
		},
		RSSFeedURL: item.RSSFeedURL,
	}
}

func SerializableItemsFromQueueItemByUserIDRows(rows []query.GetItemsWithURLByUserIDRow) []SerializableItem {
	return lop.Map(rows, func(row query.GetItemsWithURLByUserIDRow, _ int) SerializableItem {
		return SerializableItemFromQueueItemByUserIDRow(row)
	})
}

func SerializableItemFromQueueItemByUserIDRow(item query.GetItemsWithURLByUserIDRow) SerializableItem {
	return SerializableItem{
		RSSItem: &query.RSSItem{
			ID:              item.ReaderID,
			UUID:            item.UUID,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
			FeedID:          item.FeedID,
			Title:           item.Title,
			RSSGuid:         item.RSSGuid,
			Link:            item.Link,
			Author:          item.Author,
			AuthorEmail:     item.AuthorEmail,
			Content:         item.Content,
			Summary:         item.Summary,
			PublishedAt:     item.PublishedAt,
			RemoteUpdatedAt: item.RemoteUpdatedAt,
			ReaderID:        item.ReaderID,
		},
		RSSFeedURL: item.RSSFeedURL,
	}
}

func SerializableItemsFromRSSItemsAndFeedURL(items []query.RSSItem, url string) []SerializableItem {
	return lop.Map(items, func(item query.RSSItem, _ int) SerializableItem {
		return SerializableItemFromRSSItemAndFeedURL(item, url)
	})
}

func SerializableItemFromRSSItemAndFeedURL(item query.RSSItem, url string) SerializableItem {
	return SerializableItem{
		RSSItem:    &item,
		RSSFeedURL: url,
	}
}
