-- +goose Up
-- +goose StatementBegin
create schema queue;
create table queue.items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint not null references identity.users(id),
  read boolean not null default false,
  starred boolean not null default false,

  rss_item_id bigint references rss.items(id)
);

create index queue_items_index on queue.items (user_id);

create unique index queue_items_user_id_rss_item_id_key
  on queue.items(user_id, rss_item_id);

create index queue_items_user_id_read
  on queue.items (user_id, read);

create index queue_items_user_id_starred
  on queue.items (user_id, starred);

create index queue_items_rss_item_id_index
  on queue.items (created_at desc);

create trigger queue_items_touch_updated_at
  before update on queue.items for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table queue.items;
drop schema queue;
-- +goose StatementEnd
