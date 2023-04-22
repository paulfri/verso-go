-- +goose Up
-- +goose StatementBegin
create table content.queue_items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint not null references identity.users(id),
  rss_item_id bigint references content.rss_items(id),
  is_read boolean not null default false
);

create index content_queue_items_index
  on content.queue_items (user_id);

create unique index content_queue_items_user_id_rss_item_id_key
  on content.queue_items(user_id, rss_item_id);

create index content_queue_items_rss_item_id_index
  on content.queue_items (created_at desc);

create trigger content_queue_items_touch_updated_at
  before update on content.queue_items for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table content.queue_items;
-- +goose StatementEnd
