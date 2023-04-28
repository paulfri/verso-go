-- +goose NO TRANSACTION
-- +goose Up
create schema queue;
create table queue.items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  user_id bigint not null references identity.users(id) on delete cascade,
  rss_item_id bigint references rss.items(id) on delete cascade,

  -- maybe: read_at (recently read), read_version (updated items)
  read boolean not null default false,
  -- maybe: starred_at
  starred boolean not null default false
);

create index concurrently if not exists queue_items_user_id_fkey
  on queue.items (user_id);

create unique index queue_items_user_id_rss_item_id_key
  on queue.items(user_id, rss_item_id);

create index concurrently if not exists queue_items_user_id_read
  on queue.items (user_id, read);

create index concurrently if not exists queue_items_user_id_starred
  on queue.items (user_id, starred);

create index concurrently if not exists queue_items_rss_item_id_index
  on queue.items (created_at desc);

create trigger queue_items_touch_updated_at
  before update on queue.items for each row
  execute procedure touch_updated_at();

-- +goose Down
drop table queue.items;
drop schema queue;
