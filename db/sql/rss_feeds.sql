-- name: FindRssFeed :one
select * from rss.feeds
where id = $1 limit 1;

-- name: FindRssFeedByUuid :one
select * from rss.feeds
where uuid = $1 limit 1;

-- name: FindRssFeedByUrl :one
select * from rss.feeds
where url = $1 limit 1;

-- name: FindOrCreateRssFeed :one
with inserted as (
  insert into rss.feeds (
    title,
    url
  ) select $1, $2 where not exists (
    select 1 from rss.feeds where url = $2
  ) returning *
) select * from inserted
  union
  select * from rss.feeds where url = $2;

-- name: CreateRssFeed :one
insert into rss.feeds (title, url) values ($1, $2) returning *;
