// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: queue_items.sql

package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createQueueItem = `-- name: CreateQueueItem :one
insert into queue.items (user_id, rss_item_id)
values ($1, $2)
on conflict do nothing
returning id, uuid, created_at, updated_at, user_id, unread, starred, rss_item_id
`

type CreateQueueItemParams struct {
	UserID    int64         `json:"user_id"`
	RSSItemID sql.NullInt64 `json:"rss_item_id"`
}

func (q *Queries) CreateQueueItem(ctx context.Context, arg CreateQueueItemParams) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, createQueueItem, arg.UserID, arg.RSSItemID)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Unread,
		&i.Starred,
		&i.RSSItemID,
	)
	return i, err
}

const getItemsByUserID = `-- name: GetItemsByUserID :many

select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
order by ri.published_at desc
limit $2
`

type GetItemsByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
}

type GetItemsByUserIDRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

// TODO: Lots of repetition in these queries which is leading to a lot of
// boilerplate in Go. Can we do better?
func (q *Queries) GetItemsByUserID(ctx context.Context, arg GetItemsByUserIDParams) ([]GetItemsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getItemsByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetItemsByUserIDRow
	for rows.Next() {
		var i GetItemsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getItemsWithContentDataByReaderIDs = `-- name: GetItemsWithContentDataByReaderIDs :many
select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where
  ri.reader_id = any($1::text[]) -- TODO: Can't name arg due to bug in sqlc.
  and qi.user_id = $2
order by ri.published_at desc
`

type GetItemsWithContentDataByReaderIDsParams struct {
	Column1 []string `json:"column_1"`
	UserID  int64    `json:"user_id"`
}

type GetItemsWithContentDataByReaderIDsRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

func (q *Queries) GetItemsWithContentDataByReaderIDs(ctx context.Context, arg GetItemsWithContentDataByReaderIDsParams) ([]GetItemsWithContentDataByReaderIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getItemsWithContentDataByReaderIDs, pq.Array(arg.Column1), arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetItemsWithContentDataByReaderIDsRow
	for rows.Next() {
		var i GetItemsWithContentDataByReaderIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getItemsWithURLByUserID = `-- name: GetItemsWithURLByUserID :many
select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
order by ri.published_at desc
limit $2
`

type GetItemsWithURLByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
}

type GetItemsWithURLByUserIDRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

func (q *Queries) GetItemsWithURLByUserID(ctx context.Context, arg GetItemsWithURLByUserIDParams) ([]GetItemsWithURLByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getItemsWithURLByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetItemsWithURLByUserIDRow
	for rows.Next() {
		var i GetItemsWithURLByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReadItemsByUserID = `-- name: GetReadItemsByUserID :many
select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.unread = false
order by ri.published_at desc
limit $2
`

type GetReadItemsByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
}

type GetReadItemsByUserIDRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

func (q *Queries) GetReadItemsByUserID(ctx context.Context, arg GetReadItemsByUserIDParams) ([]GetReadItemsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getReadItemsByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReadItemsByUserIDRow
	for rows.Next() {
		var i GetReadItemsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStarredItemsByUserID = `-- name: GetStarredItemsByUserID :many
select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.starred = true
order by ri.published_at desc
limit $2
`

type GetStarredItemsByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
}

type GetStarredItemsByUserIDRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

func (q *Queries) GetStarredItemsByUserID(ctx context.Context, arg GetStarredItemsByUserIDParams) ([]GetStarredItemsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getStarredItemsByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStarredItemsByUserIDRow
	for rows.Next() {
		var i GetStarredItemsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnreadCountsByUserID = `-- name: GetUnreadCountsByUserID :many
select rf.id, rf.url, count(*), max(ri.published_at) as newest from queue.items qi
  join rss.items ri on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where user_id = $1
  and unread = true
group by rf.id
`

type GetUnreadCountsByUserIDRow struct {
	ID     int64       `json:"id"`
	URL    string      `json:"url"`
	Count  int64       `json:"count"`
	Newest interface{} `json:"newest"`
}

func (q *Queries) GetUnreadCountsByUserID(ctx context.Context, userID int64) ([]GetUnreadCountsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getUnreadCountsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUnreadCountsByUserIDRow
	for rows.Next() {
		var i GetUnreadCountsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.URL,
			&i.Count,
			&i.Newest,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnreadItemsByUserID = `-- name: GetUnreadItemsByUserID :many
select ri.id, ri.uuid, ri.created_at, ri.updated_at, ri.feed_id, ri.rss_guid, ri.title, ri.link, ri.author, ri.author_email, ri.content, ri.summary, ri.published_at, ri.remote_updated_at, ri.reader_id, rf.url as rss_feed_url, qi.user_id, qi.unread, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.unread = true
order by ri.published_at desc
limit $2
`

type GetUnreadItemsByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
}

type GetUnreadItemsByUserIDRow struct {
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
	RSSFeedURL      string         `json:"rss_feed_url"`
	UserID          int64          `json:"user_id"`
	Unread          bool           `json:"unread"`
	Starred         bool           `json:"starred"`
}

func (q *Queries) GetUnreadItemsByUserID(ctx context.Context, arg GetUnreadItemsByUserIDParams) ([]GetUnreadItemsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getUnreadItemsByUserID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUnreadItemsByUserIDRow
	for rows.Next() {
		var i GetUnreadItemsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.RSSGuid,
			&i.Title,
			&i.Link,
			&i.Author,
			&i.AuthorEmail,
			&i.Content,
			&i.Summary,
			&i.PublishedAt,
			&i.RemoteUpdatedAt,
			&i.ReaderID,
			&i.RSSFeedURL,
			&i.UserID,
			&i.Unread,
			&i.Starred,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markAllQueueItemsAsRead = `-- name: MarkAllQueueItemsAsRead :exec
update queue.items qi
  set unread = false
  where exists (
    select ri.id, ri.uuid, ri.created_at, ri.updated_at, feed_id, rss_guid, ri.title, link, author, author_email, content, summary, published_at, remote_updated_at, reader_id, rf.id, rf.uuid, rf.created_at, rf.updated_at, rf.title, url, active, last_crawled_at from rss.items ri
      join rss.feeds rf on rf.id = ri.feed_id
    where
      qi.rss_item_id = ri.id
      and qi.user_id = $1
      and rf.url = coalesce(nullif($2,''), rf.url)
      and ri.published_at <= coalesce($3, now())
  )
`

type MarkAllQueueItemsAsReadParams struct {
	UserID          int64        `json:"user_id"`
	RSSFeedURL      interface{}  `json:"rss_feed_url"`
	PublishedBefore sql.NullTime `json:"published_before"`
}

func (q *Queries) MarkAllQueueItemsAsRead(ctx context.Context, arg MarkAllQueueItemsAsReadParams) error {
	_, err := q.db.ExecContext(ctx, markAllQueueItemsAsRead, arg.UserID, arg.RSSFeedURL, arg.PublishedBefore)
	return err
}

const updateQueueItemReadState = `-- name: UpdateQueueItemReadState :one
update queue.items qi
  set unread = $1
from rss.items ri
where
  qi.rss_item_id = ri.id
  and ri.reader_id = $2
  and qi.user_id = $3
returning qi.id, qi.uuid, qi.created_at, qi.updated_at, qi.user_id, qi.unread, qi.starred, qi.rss_item_id
`

type UpdateQueueItemReadStateParams struct {
	Unread   bool   `json:"unread"`
	ReaderID string `json:"reader_id"`
	UserID   int64  `json:"user_id"`
}

func (q *Queries) UpdateQueueItemReadState(ctx context.Context, arg UpdateQueueItemReadStateParams) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, updateQueueItemReadState, arg.Unread, arg.ReaderID, arg.UserID)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Unread,
		&i.Starred,
		&i.RSSItemID,
	)
	return i, err
}

const updateQueueItemStarredState = `-- name: UpdateQueueItemStarredState :one
update queue.items qi
  set starred = $1
from rss.items ri
where
  qi.rss_item_id = ri.id
  and ri.reader_id = $2
  and qi.user_id = $3
returning qi.id, qi.uuid, qi.created_at, qi.updated_at, qi.user_id, qi.unread, qi.starred, qi.rss_item_id
`

type UpdateQueueItemStarredStateParams struct {
	Starred  bool   `json:"starred"`
	ReaderID string `json:"reader_id"`
	UserID   int64  `json:"user_id"`
}

func (q *Queries) UpdateQueueItemStarredState(ctx context.Context, arg UpdateQueueItemStarredStateParams) (QueueItem, error) {
	row := q.db.QueryRowContext(ctx, updateQueueItemStarredState, arg.Starred, arg.ReaderID, arg.UserID)
	var i QueueItem
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Unread,
		&i.Starred,
		&i.RSSItemID,
	)
	return i, err
}
