// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package query

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type IdentityReaderToken struct {
	ID         int64        `json:"id"`
	UUID       uuid.UUID    `json:"uuid"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	UserID     int64        `json:"user_id"`
	Identifier string       `json:"identifier"`
	RevokedAt  sql.NullTime `json:"revoked_at"`
}

type IdentityUser struct {
	ID        int64          `json:"id"`
	UUID      uuid.UUID      `json:"uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	Password  sql.NullString `json:"password"`
	Superuser bool           `json:"superuser"`
}

type QueueItem struct {
	ID             int64     `json:"id"`
	UUID           uuid.UUID `json:"uuid"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         int64     `json:"user_id"`
	RSSItemID      int64     `json:"rss_item_id"`
	SubscriptionID int64     `json:"subscription_id"`
	Read           bool      `json:"read"`
	Starred        bool      `json:"starred"`
}

type RSSFeed struct {
	ID            int64        `json:"id"`
	UUID          uuid.UUID    `json:"uuid"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Title         string       `json:"title"`
	URL           string       `json:"url"`
	Active        bool         `json:"active"`
	LastCrawledAt sql.NullTime `json:"last_crawled_at"`
}

type RSSItem struct {
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
}

type RSSItemVersion struct {
	ID              int64         `json:"id"`
	UUID            uuid.UUID     `json:"uuid"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ItemID          sql.NullInt64 `json:"item_id"`
	Title           string        `json:"title"`
	Link            string        `json:"link"`
	Content         string        `json:"content"`
	PublishedAt     sql.NullTime  `json:"published_at"`
	RemoteUpdatedAt sql.NullTime  `json:"remote_updated_at"`
}

type RSSSubscription struct {
	ID          int64          `json:"id"`
	UUID        uuid.UUID      `json:"uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UserID      int64          `json:"user_id"`
	FeedID      int64          `json:"feed_id"`
	CustomTitle sql.NullString `json:"custom_title"`
}

type TaxonomyRssFeedTag struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TagID     int64     `json:"tag_id"`
	RSSFeedID int64     `json:"rss_feed_id"`
}

type TaxonomyTag struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
}
