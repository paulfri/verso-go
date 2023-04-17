-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1 LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM feeds;
