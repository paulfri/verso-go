package serialize

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	ID              int64          `json:"id"`
	UUID            uuid.UUID      `json:"uuid"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	FeedID          int64          `json:"feed_id"`
	RSSGuid         string         `json:"rss_guid"`
	Title           string         `json:"title"`
	Link            string         `json:"link"`
	Author          sql.NullString `json:"author"`
	AuthorEmail     sql.NullString `json:"author_email"`
	Content         string         `json:"content"`
	Summary         sql.NullString `json:"summary"`
	PublishedAt     sql.NullTime   `json:"published_at"`
	RemoteUpdatedAt sql.NullTime   `json:"remote_updated_at"`
	ReaderID        string         `json:"reader_id"`
	UserID          int64          `json:"user_id"`
	RSSFeedURL      string         `json:"rss_feed_url"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

// Settle down, this one is good.
//
// In sqlc, SELECT items.* will return []Item, but SELECT items.*, feeds.col
// will return a struct specific to that query. So if you have a bunch of
// queries that select something like that, each query will have its own struct
// result, even though they are identical.
//
// Compounding that, Go generics and interfaces don't support struct fields. So
// if you have a bunch of queries that all return the same fields, you can't use
// them interoperably in Go functions. Support for struct fields in generics
// might be coming eventually? https://github.com/golang/go/issues/48522
//
// This would also resolve if sqlc added the ability to name or alias query row
// structs. https://github.com/kyleconroy/sqlc/discussions/1183
//
// The workaround is to marshal the structs to JSON and back, since they have
// identical JSON serialization tags. This is a huge hack and probably very
// slow?
func QueryRowsToSerializableItems(row interface{}) []SerializableItem {
	asJSON, err := json.Marshal(row)
	if err != nil {
		panic(err)
	}

	var items []SerializableItem
	err = json.Unmarshal(asJSON, &items)
	if err != nil {
		panic(err)
	}

	return items
}
