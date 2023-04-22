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
