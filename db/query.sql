-- name: GetRssFeedById :one
select * from content.rss_feeds
where id = $1 limit 1;

-- name: GetRssFeedByUuid :one
select * from content.rss_feeds
where uuid = $1 limit 1;

-- name: ListRSSFeeds :many
select * from content.rss_feeds;

-- name: CreateRssFeed :one
insert into content.rss_feeds (title, url) values ($1, $2) returning *;

-- name: CreateItem :one
insert into content.rss_items (
  rss_feed_id,
  rss_guid,
  title,
  link,
  content,
  published_at,
  remote_updated_at
) values ($1, $2, $3, $4, $5, $6, $7) returning *;
