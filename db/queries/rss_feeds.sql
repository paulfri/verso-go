-- name: FindRssFeed :one
select * from content.rss_feeds
where id = $1 limit 1;

-- name: FindRssFeedByUuid :one
select * from content.rss_feeds
where uuid = $1 limit 1;

-- name: FindOrCreateRssFeed :one
with inserted as (
  insert into content.rss_feeds (
    title,
    url
  ) select $1, $2 where not exists (
    select 1 from content.rss_feeds where url = $2
  ) returning *
) select * from inserted union select * from content.rss_feeds where url = $2;

-- name: CreateRssFeed :one
insert into content.rss_feeds (title, url) values ($1, $2) returning *;

-- name: CreateItem :one
insert into content.rss_items as i (
  rss_feed_id,
  rss_guid,
  title,
  link,
  content,
  published_at,
  remote_updated_at
) values ($1, $2, $3, $4, $5, $6, $7)
  on conflict (rss_feed_id, rss_guid) do update
  set
    title = excluded.title, 
    link = excluded.link, 
    content = excluded.content,
    published_at = excluded.published_at,
    remote_updated_at = excluded.remote_updated_at
  -- only perform insert if these values are changed, in order to avoid the
  -- rss_item_versions_insert_on_item_update trigger on updates to this table
  where
    i.title is distinct from excluded.title or 
    i.link is distinct from excluded.link or 
    i.content is distinct from excluded.content or
    i.published_at is distinct from excluded.published_at or
    i.remote_updated_at is distinct from excluded.remote_updated_at
  returning *;

-- name: CreateSubscription :one
insert into content.rss_subscriptions (
  user_id,
  rss_feed_id
) select $1, $2 where not exists (
  select 1 from content.rss_subscriptions where user_id = $1 and rss_feed_id = $2
) returning *;
