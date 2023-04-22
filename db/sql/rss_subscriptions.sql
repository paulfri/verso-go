-- name: CreateSubscription :one
with inserted as (
  insert into content.rss_subscriptions (
    user_id,
    rss_feed_id
  ) select $1, $2 where not exists (
    select 1 from content.rss_subscriptions where user_id = $1 and rss_feed_id = $2
  ) returning *
)
select * from inserted
  union
  select * from content.rss_subscriptions where user_id = $1 and rss_feed_id = $2;

-- name: GetSubscribersByFeedId :many
select * from content.rss_subscriptions
  where rss_subscriptions.rss_feed_id = $1;
