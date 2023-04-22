-- name: GetQueueItemsByUserID :many
select * from rss.items ri
  join queue.items qi on qi.rss_item_id = ri.id
  where qi.user_id = $1 
  order by ri.published_at desc
  limit $2;

-- name: CreateQueueItem :one
insert into queue.items (user_id, rss_item_id)
  values ($1, $2)
  on conflict do nothing
  returning *;
