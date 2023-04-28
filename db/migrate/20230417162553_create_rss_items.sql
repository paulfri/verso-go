-- +goose NO TRANSACTION
-- +goose Up
create table rss.items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  feed_id bigint not null references rss.feeds(id) on delete cascade,

  rss_guid text not null,
  title text not null,
  link text not null,
  author text,
  author_email text,
  content text not null,
  summary text,
  published_at timestamptz,
  remote_updated_at timestamptz,
  reader_id text unique not null constraint reader_id_length_check
    check (char_length(reader_id) = 16)
);

create unique index concurrently if not exists rss_items_rss_feed_id_rss_guid_key
  on rss.items (feed_id, rss_guid);

create index concurrently if not exists rss_items_published_at_idx
  on rss.items (published_at desc);

create trigger rss_items_generate_reader_id
  before insert on rss.items for each row
  execute procedure generate_reader_id();

create trigger rss_items_touch_updated_at
  before update on rss.items for each row
  execute procedure touch_updated_at();

-- +goose Down
drop function generate_reader_id();
drop table rss.items;
