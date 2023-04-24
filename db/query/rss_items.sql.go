// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: rss_items.sql

package query

import (
	"context"
	"database/sql"
)

const createRSSItem = `-- name: CreateRSSItem :one
insert into rss.items as items (
  feed_id,
  rss_guid,
  title,
  link,
  author,
  author_email,
  content,
  published_at,
  remote_updated_at
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  on conflict (feed_id, rss_guid) do update
  set
    title = excluded.title, 
    link = excluded.link, 
    content = excluded.content,
    published_at = excluded.published_at,
    remote_updated_at = excluded.remote_updated_at
  -- only perform insert if these values are changed, in order to avoid the
  -- rss_item_versions_insert_on_item_update trigger on updates to this table
  where
    items.title is distinct from excluded.title or 
    items.link is distinct from excluded.link or 
    items.content is distinct from excluded.content or
    items.published_at is distinct from excluded.published_at or
    items.remote_updated_at is distinct from excluded.remote_updated_at
  returning id, uuid, created_at, updated_at, feed_id, rss_guid, title, link, author, author_email, content, summary, published_at, remote_updated_at, reader_id
`

type CreateRSSItemParams struct {
	FeedID          int64          `json:"feed_id"`
	RSSGuid         string         `json:"rss_guid"`
	Title           string         `json:"title"`
	Link            string         `json:"link"`
	Author          sql.NullString `json:"author"`
	AuthorEmail     sql.NullString `json:"author_email"`
	Content         string         `json:"content"`
	PublishedAt     sql.NullTime   `json:"published_at"`
	RemoteUpdatedAt sql.NullTime   `json:"remote_updated_at"`
}

func (q *Queries) CreateRSSItem(ctx context.Context, arg CreateRSSItemParams) (RSSItem, error) {
	row := q.db.QueryRowContext(ctx, createRSSItem,
		arg.FeedID,
		arg.RSSGuid,
		arg.Title,
		arg.Link,
		arg.Author,
		arg.AuthorEmail,
		arg.Content,
		arg.PublishedAt,
		arg.RemoteUpdatedAt,
	)
	var i RSSItem
	err := row.Scan(
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
	)
	return i, err
}

const getRecentItemsByRSSFeedID = `-- name: GetRecentItemsByRSSFeedID :many
select id, uuid, created_at, updated_at, feed_id, rss_guid, title, link, author, author_email, content, summary, published_at, remote_updated_at, reader_id from rss.items
  where items.feed_id = $1
  order by items.published_at desc
  limit $2
`

type GetRecentItemsByRSSFeedIDParams struct {
	FeedID int64 `json:"feed_id"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) GetRecentItemsByRSSFeedID(ctx context.Context, arg GetRecentItemsByRSSFeedIDParams) ([]RSSItem, error) {
	rows, err := q.db.QueryContext(ctx, getRecentItemsByRSSFeedID, arg.FeedID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RSSItem
	for rows.Next() {
		var i RSSItem
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
