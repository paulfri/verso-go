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
