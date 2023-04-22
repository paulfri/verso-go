-- +goose Up
-- +goose StatementBegin
create table content.rss_items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  rss_feed_id bigint not null references content.rss_feeds(id) on delete cascade,

  rss_guid text not null,
  title text not null,
  link text not null,
  content text not null,
  published_at timestamp with time zone,
  remote_updated_at timestamp with time zone
);

create trigger content_rss_items_touch_updated_at
  before update on content.rss_items for each row
  execute procedure touch_updated_at();

create unique index content_rss_items_rss_feed_id_rss_guid_key
  on content.rss_items (rss_feed_id, rss_guid);
create index content_rss_items_published_at_index
  on content.rss_items (published_at desc);
create index content_rss_items_updated_at_index
  on content.rss_items (updated_at desc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table content.rss_items;
-- +goose StatementEnd
