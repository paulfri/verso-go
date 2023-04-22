-- name: GetQueueItemsByUserId :many
select * from content.rss_items ri
  join content.queue_items on queue_items.rss_item_id = ri.id
  where queue_items.user_id = $1 
  order by ri.id desc
  limit $2;

-- name: CreateQueueItem :one
insert into content.queue_items (user_id, rss_item_id)
  values ($1, $2)
  on conflict do nothing
  returning *;
