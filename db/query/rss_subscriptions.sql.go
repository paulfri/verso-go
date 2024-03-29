// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: rss_subscriptions.sql

package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRSSSubscription = `-- name: CreateRSSSubscription :one
with inserted as (
  insert into rss.subscriptions (
    user_id,
    feed_id
  ) select $1, $2 where not exists (
    select 1 from rss.subscriptions where user_id = $1 and feed_id = $2
  ) returning id, uuid, created_at, updated_at, user_id, feed_id, custom_title
)
select id, uuid, created_at, updated_at, user_id, feed_id, custom_title from inserted
  union
  select id, uuid, created_at, updated_at, user_id, feed_id, custom_title from rss.subscriptions where user_id = $1 and feed_id = $2
`

type CreateRSSSubscriptionParams struct {
	UserID int64 `json:"user_id"`
	FeedID int64 `json:"feed_id"`
}

type CreateRSSSubscriptionRow struct {
	ID          int64          `json:"id"`
	UUID        uuid.UUID      `json:"uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UserID      int64          `json:"user_id"`
	FeedID      int64          `json:"feed_id"`
	CustomTitle sql.NullString `json:"custom_title"`
}

func (q *Queries) CreateRSSSubscription(ctx context.Context, arg CreateRSSSubscriptionParams) (CreateRSSSubscriptionRow, error) {
	row := q.db.QueryRowContext(ctx, createRSSSubscription, arg.UserID, arg.FeedID)
	var i CreateRSSSubscriptionRow
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.CustomTitle,
	)
	return i, err
}

const deleteSubscriptionByRSSFeedURLAndUserID = `-- name: DeleteSubscriptionByRSSFeedURLAndUserID :exec
delete from rss.subscriptions s
using rss.feeds f
where s.feed_id = f.id
  and f.url = $2
  and s.user_id = $1
`

type DeleteSubscriptionByRSSFeedURLAndUserIDParams struct {
	UserID     int64  `json:"user_id"`
	RSSFeedURL string `json:"rss_feed_url"`
}

func (q *Queries) DeleteSubscriptionByRSSFeedURLAndUserID(ctx context.Context, arg DeleteSubscriptionByRSSFeedURLAndUserIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteSubscriptionByRSSFeedURLAndUserID, arg.UserID, arg.RSSFeedURL)
	return err
}

const getSubscriptionByRSSFeedIDAndUserID = `-- name: GetSubscriptionByRSSFeedIDAndUserID :one
select id, uuid, created_at, updated_at, user_id, feed_id, custom_title from rss.subscriptions s
  where s.feed_id = $1 and s.user_id = $2
`

type GetSubscriptionByRSSFeedIDAndUserIDParams struct {
	FeedID int64 `json:"feed_id"`
	UserID int64 `json:"user_id"`
}

func (q *Queries) GetSubscriptionByRSSFeedIDAndUserID(ctx context.Context, arg GetSubscriptionByRSSFeedIDAndUserIDParams) (RSSSubscription, error) {
	row := q.db.QueryRowContext(ctx, getSubscriptionByRSSFeedIDAndUserID, arg.FeedID, arg.UserID)
	var i RSSSubscription
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.CustomTitle,
	)
	return i, err
}

const getSubscriptionsByRSSFeedID = `-- name: GetSubscriptionsByRSSFeedID :many
select id, uuid, created_at, updated_at, user_id, feed_id, custom_title from rss.subscriptions s
  where s.feed_id = $1
`

func (q *Queries) GetSubscriptionsByRSSFeedID(ctx context.Context, feedID int64) ([]RSSSubscription, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriptionsByRSSFeedID, feedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RSSSubscription
	for rows.Next() {
		var i RSSSubscription
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.CustomTitle,
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

const getSubscriptionsByUserID = `-- name: GetSubscriptionsByUserID :many
select
  s.id, s.uuid, s.created_at, s.updated_at, s.user_id, s.feed_id, s.custom_title,
  t.name,
  f.url as rss_feed_url,
  coalesce(s.custom_title, f.title) as title
from rss.subscriptions s
  join rss.feeds f on f.id = s.feed_id
  left outer join taxonomy.rss_feed_tags ft on f.id = ft.rss_feed_id
  left outer join taxonomy.tags t on t.id = ft.tag_id
  where s.user_id = $1
`

type GetSubscriptionsByUserIDRow struct {
	ID          int64          `json:"id"`
	UUID        uuid.UUID      `json:"uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UserID      int64          `json:"user_id"`
	FeedID      int64          `json:"feed_id"`
	CustomTitle sql.NullString `json:"custom_title"`
	Name        sql.NullString `json:"name"`
	RSSFeedURL  string         `json:"rss_feed_url"`
	Title       string         `json:"title"`
}

func (q *Queries) GetSubscriptionsByUserID(ctx context.Context, userID int64) ([]GetSubscriptionsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getSubscriptionsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSubscriptionsByUserIDRow
	for rows.Next() {
		var i GetSubscriptionsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.CustomTitle,
			&i.Name,
			&i.RSSFeedURL,
			&i.Title,
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
