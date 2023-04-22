// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: rss_subscriptions.sql

package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createSubscription = `-- name: CreateSubscription :one
with inserted as (
  insert into content.rss_subscriptions (
    user_id,
    rss_feed_id
  ) select $1, $2 where not exists (
    select 1 from content.rss_subscriptions where user_id = $1 and rss_feed_id = $2
  ) returning id, uuid, created_at, updated_at, user_id, rss_feed_id, custom_title
)
select id, uuid, created_at, updated_at, user_id, rss_feed_id, custom_title from inserted
  union
  select id, uuid, created_at, updated_at, user_id, rss_feed_id, custom_title from content.rss_subscriptions where user_id = $1 and rss_feed_id = $2
`

type CreateSubscriptionParams struct {
	UserID    int64 `json:"user_id"`
	RssFeedID int64 `json:"rss_feed_id"`
}

type CreateSubscriptionRow struct {
	ID          int64          `json:"id"`
	Uuid        uuid.UUID      `json:"uuid"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	UserID      int64          `json:"user_id"`
	RssFeedID   int64          `json:"rss_feed_id"`
	CustomTitle sql.NullString `json:"custom_title"`
}

func (q *Queries) CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) (CreateSubscriptionRow, error) {
	row := q.db.QueryRowContext(ctx, createSubscription, arg.UserID, arg.RssFeedID)
	var i CreateSubscriptionRow
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.RssFeedID,
		&i.CustomTitle,
	)
	return i, err
}
