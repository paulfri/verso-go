-- +goose Up
-- +goose StatementBegin
create table rss.items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  feed_id bigint not null references rss.feeds(id) on delete cascade,

  rss_guid text not null,
  title text not null,
  link text not null,
  author text,
  author_email text,
  content text not null,
  summary text,
  published_at timestamp with time zone,
  remote_updated_at timestamp with time zone
);

create trigger rss_items_touch_updated_at
  before update on rss.items for each row
  execute procedure touch_updated_at();

create unique index rss_items_rss_feed_id_rss_guid_key
  on rss.items (feed_id, rss_guid);
create index rss_items_published_at_index
  on rss.items (published_at desc);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rss.items;
-- +goose StatementEnd
