-- name: GetQueueItemsByUserID :many
select
    qi.id,
    qi.created_at,
    qi.unread,
    ri.feed_id,
    ri.title,
    ri.rss_guid,
    ri.link,
    ri.author,
    ri.content,
    ri.summary,
    ri.published_at,
    ri.remote_updated_at
  from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  where qi.user_id = $1 
  order by ri.published_at desc
  limit $2;

-- name: CreateQueueItem :one
insert into queue.items (user_id, rss_item_id)
  values ($1, $2)
  on conflict do nothing
  returning *;

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
