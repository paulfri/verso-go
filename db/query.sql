-- name: GetFeedById :one
SELECT * FROM feeds
WHERE id = $1 LIMIT 1;

-- name: GetFeedByUuid :one
SELECT * FROM feeds
WHERE uuid = $1 LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: CreateFeed :one
INSERT INTO feeds (title, url) VALUES ($1, $2)
RETURNING *;

-- name: CreateItem :one
INSERT INTO items (
  feed_id,
  remote_id,
  title,
  link,
  content,
  published_at,
  remote_updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;
