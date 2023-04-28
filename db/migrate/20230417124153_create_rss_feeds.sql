-- +goose Up
-- +goose StatementBegin
create schema rss;

create table rss.feeds (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  title text not null,
  url text not null unique,
  active boolean not null default true,
  last_crawled_at timestamp
);

create index rss_feeds_url_idx
  on rss.feeds (url);

create index rss_feeds_last_crawled_at_idx
  on rss.feeds (last_crawled_at);

create trigger rss_feeds_touch_updated_at
  before update on rss.feeds for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rss.feeds;
drop schema rss;
-- +goose StatementEnd
