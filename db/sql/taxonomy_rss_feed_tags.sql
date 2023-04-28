-- name: CreateRSSFeedTag :one
insert into taxonomy.rss_feed_tags (
  tag_id,
  rss_feed_id
)
values ($1, $2)
returning *;
