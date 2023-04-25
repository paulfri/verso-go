-- name: GetRecentItemsByRSSFeedID :many
select ri.* from rss.items ri
where ri.feed_id = $1
order by ri.published_at desc
limit $2;

-- name: GetRecentItemsByRSSFeedURL :many
select ri.* from rss.items ri
  join rss.feeds rf on rf.id = ri.feed_id
where rf.url = $1
order by ri.published_at desc
limit $2;

-- name: CreateRSSItem :one
insert into rss.items as items (
  feed_id,
  rss_guid,
  title,
  link,
  author,
  author_email,
  content,
  published_at,
  remote_updated_at
) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  on conflict (feed_id, rss_guid) do update
  set
    title = excluded.title, 
    link = excluded.link, 
    content = excluded.content,
    published_at = excluded.published_at,
    remote_updated_at = excluded.remote_updated_at
  -- only perform insert if these values are changed, in order to avoid the
  -- rss_item_versions_insert_on_item_update trigger on updates to this table
  where
    items.title is distinct from excluded.title or 
    items.link is distinct from excluded.link or 
    items.content is distinct from excluded.content or
    items.published_at is distinct from excluded.published_at or
    items.remote_updated_at is distinct from excluded.remote_updated_at
  returning *;
