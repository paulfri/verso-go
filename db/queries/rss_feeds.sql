-- name: FindRssFeed :one
select * from content.rss_feeds
where id = $1 limit 1;

-- name: FindRssFeedByUuid :one
select * from content.rss_feeds
where uuid = $1 limit 1;

-- name: FindRssFeedByUrl :one
select * from content.rss_feeds
where url = $1 limit 1;

-- name: FindOrCreateRssFeed :one
with inserted as (
  insert into content.rss_feeds (
    title,
    url
  ) select $1, $2 where not exists (
    select 1 from content.rss_feeds where url = $2
  ) returning *
) select * from inserted
  union
  select * from content.rss_feeds where url = $2;

-- name: CreateRssFeed :one
insert into content.rss_feeds (title, url) values ($1, $2) returning *;
