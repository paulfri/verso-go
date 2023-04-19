// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: query.sql

package query

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createItem = `-- name: CreateItem :one
insert into content.rss_items as i (
  rss_feed_id,
  rss_guid,
  title,
  link,
  content,
  published_at,
  remote_updated_at
) values ($1, $2, $3, $4, $5, $6, $7)
  on conflict (rss_feed_id, rss_guid) do update
  set
    title = excluded.title, 
    link = excluded.link, 
    content = excluded.content,
    published_at = excluded.published_at,
    remote_updated_at = excluded.remote_updated_at
  -- only perform insert if these values are changed, in order to avoid the
  -- rss_item_versions_insert_on_item_update trigger on updates to this table
  where
    i.title is distinct from excluded.title or 
    i.link is distinct from excluded.link or 
    i.content is distinct from excluded.content or
    i.published_at is distinct from excluded.published_at or
    i.remote_updated_at is distinct from excluded.remote_updated_at
  returning id, uuid, created_at, updated_at, rss_feed_id, rss_guid, title, link, content, published_at, remote_updated_at
`

type CreateItemParams struct {
	RssFeedID       int64        `json:"rss_feed_id"`
	RssGuid         string       `json:"rss_guid"`
	Title           string       `json:"title"`
	Link            string       `json:"link"`
	Content         string       `json:"content"`
	PublishedAt     sql.NullTime `json:"published_at"`
	RemoteUpdatedAt sql.NullTime `json:"remote_updated_at"`
}

func (q *Queries) CreateItem(ctx context.Context, arg CreateItemParams) (ContentRssItem, error) {
	row := q.db.QueryRowContext(ctx, createItem,
		arg.RssFeedID,
		arg.RssGuid,
		arg.Title,
		arg.Link,
		arg.Content,
		arg.PublishedAt,
		arg.RemoteUpdatedAt,
	)
	var i ContentRssItem
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RssFeedID,
		&i.RssGuid,
		&i.Title,
		&i.Link,
		&i.Content,
		&i.PublishedAt,
		&i.RemoteUpdatedAt,
	)
	return i, err
}

const createRssFeed = `-- name: CreateRssFeed :one
insert into content.rss_feeds (title, url) values ($1, $2) returning id, uuid, created_at, updated_at, title, url, active, last_crawled_at
`

type CreateRssFeedParams struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (q *Queries) CreateRssFeed(ctx context.Context, arg CreateRssFeedParams) (ContentRssFeed, error) {
	row := q.db.QueryRowContext(ctx, createRssFeed, arg.Title, arg.Url)
	var i ContentRssFeed
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Active,
		&i.LastCrawledAt,
	)
	return i, err
}

const getRssFeedById = `-- name: GetRssFeedById :one
select id, uuid, created_at, updated_at, title, url, active, last_crawled_at from content.rss_feeds
where id = $1 limit 1
`

func (q *Queries) GetRssFeedById(ctx context.Context, id int64) (ContentRssFeed, error) {
	row := q.db.QueryRowContext(ctx, getRssFeedById, id)
	var i ContentRssFeed
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Active,
		&i.LastCrawledAt,
	)
	return i, err
}

const getRssFeedByUuid = `-- name: GetRssFeedByUuid :one
select id, uuid, created_at, updated_at, title, url, active, last_crawled_at from content.rss_feeds
where uuid = $1 limit 1
`

func (q *Queries) GetRssFeedByUuid(ctx context.Context, uuid uuid.UUID) (ContentRssFeed, error) {
	row := q.db.QueryRowContext(ctx, getRssFeedByUuid, uuid)
	var i ContentRssFeed
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Active,
		&i.LastCrawledAt,
	)
	return i, err
}

const listRSSFeeds = `-- name: ListRSSFeeds :many
select id, uuid, created_at, updated_at, title, url, active, last_crawled_at from content.rss_feeds
`

func (q *Queries) ListRSSFeeds(ctx context.Context) ([]ContentRssFeed, error) {
	rows, err := q.db.QueryContext(ctx, listRSSFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ContentRssFeed
	for rows.Next() {
		var i ContentRssFeed
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Active,
			&i.LastCrawledAt,
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
