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
select * from rss.subscriptions s
  join rss.feeds f on f.id = s.feed_id
  where s.user_id = $1;
