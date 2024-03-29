-- TODO: Lots of repetition in these queries which is leading to a lot of
-- boilerplate in Go. Can we do better?

-- name: GetItemsByUserID :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
order by ri.published_at desc
limit $2;

-- name: GetReadItemsByUserID :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.read = true
order by ri.published_at desc
limit $2;

-- name: GetUnreadItemsByUserID :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.read = false
order by ri.published_at desc
limit $2;

-- name: GetStarredItemsByUserID :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
  and qi.starred = true
order by ri.published_at desc
limit $2;

-- name: GetItemsWithURLByUserID :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where qi.user_id = $1 
order by ri.published_at desc
limit $2;

-- name: GetItemsWithContentDataByReaderIDs :many
select ri.*, rf.url as rss_feed_url, qi.user_id, qi.read, qi.starred
from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where
  ri.reader_id = any($1::text[]) -- TODO: Can't name arg due to bug in sqlc.
  and qi.user_id = $2
order by ri.published_at desc;

-- name: GetUnreadCountsByUserID :many
select rf.id, rf.url, count(*), max(ri.published_at) as newest from queue.items qi
  join rss.items ri on qi.rss_item_id = ri.id
  join rss.feeds rf on ri.feed_id = rf.id
where user_id = $1
  and qi.read = false
group by rf.id;

-- name: CreateQueueItem :one
insert into queue.items (user_id, rss_item_id, subscription_id)
values ($1, $2, $3)
on conflict do nothing
returning *;

-- name: UpdateQueueItemReadState :one
update queue.items qi
  set read = $1
from rss.items ri
where
  qi.rss_item_id = ri.id
  and ri.reader_id = @reader_id
  and qi.user_id = @user_id
returning qi.*;

-- name: UpdateQueueItemStarredState :one
update queue.items qi
  set starred = $1
from rss.items ri
where
  qi.rss_item_id = ri.id
  and ri.reader_id = @reader_id
  and qi.user_id = @user_id
returning qi.*;

-- name: MarkAllQueueItemsAsRead :exec
update queue.items qi
  set read = true
  where exists (
    select * from rss.items ri
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
--   set read = true
--   from rss.items ri
--     join rss.feeds rf on rf.id = ri.feed_id
--   where
--     qi.rss_item_id = ri.id and qi.user_id = $1 and rf.url = $2
--   returning *;
