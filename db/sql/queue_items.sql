-- name: GetQueueItemsByUserID :many
select ri.*, rf.url as rss_feed_url
  from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
  where qi.user_id = $1 
  order by ri.published_at desc
  limit $2;

-- name: CreateQueueItem :one
insert into queue.items (user_id, rss_item_id)
  values ($1, $2)
  on conflict do nothing
  returning *;

-- name: GetQueueItemsByReaderIDs :many
select ri.*, rf.url as rss_feed_url
  from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
  where ri.id = any($1::bigint[])
  order by ri.published_at desc;

-- name: MarkAllQueueItemsAsRead :exec
update queue.items qi
  set unread = false
  where exists (
    select * 
    from rss.items ri
      join rss.feeds rf on rf.id = ri.feed_id
    where
      qi.rss_item_id = ri.id
      and qi.user_id = $1
      and rf.url = coalesce(nullif(@rss_feed_url,''), rf.url)
      and ri.published_at <= coalesce(@published_before, now())
  );
-- Unfortunately this doesn't compile in sqlc. This would be much faster than
-- the alternative (ANSI) implementation above.
-- https://github.com/kyleconroy/sqlc/issues/1100
-- update queue.items qi
--   set qi.unread = false
--   from rss.items ri
--     join rss.feeds rf on rf.id = ri.feed_id
--   where
--     qi.rss_item_id = ri.id and qi.user_id = $1 and rf.url = $2
--   returning *;
