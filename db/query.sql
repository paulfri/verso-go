-- name: GetFeed :one
SELECT * FROM feeds
WHERE uuid = $1 LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM feeds;
