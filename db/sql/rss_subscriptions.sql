-- name: CreateRSSSubscription :one
with inserted as (
  insert into rss.subscriptions (
    user_id,
    feed_id
  ) select $1, $2 where not exists (
    select 1 from rss.subscriptions where user_id = $1 and feed_id = $2
  ) returning *
)
select * from inserted
  union
  select * from rss.subscriptions where user_id = $1 and feed_id = $2;

-- name: GetSubscriptionsByRSSFeedID :many
select * from rss.subscriptions s
  where s.feed_id = $1;

-- name: GetSubscriptionByRSSFeedIDAndUserID :one
select * from rss.subscriptions s
  where s.feed_id = $1 and s.user_id = $2;

-- name: GetSubscriptionsByUserID :many
select
  s.*,
  t.name,
  f.url as rss_feed_url,
  coalesce(s.custom_title, f.title) as title
from rss.subscriptions s
  join rss.feeds f on f.id = s.feed_id
  left outer join taxonomy.rss_feed_tags ft on f.id = ft.rss_feed_id
  left outer join taxonomy.tags t on t.id = ft.tag_id
  where s.user_id = $1;
