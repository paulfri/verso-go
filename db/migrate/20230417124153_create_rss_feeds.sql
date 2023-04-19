-- +goose Up
-- +goose StatementBegin
create schema content;

create table content.rss_feeds (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  title text not null,
  url text not null,
  active boolean not null default true,
  last_crawled_at timestamp
);

create trigger feeds_touch_updated_at
  before update on content.rss_feeds for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table content.rss_feeds;
drop schema content;
-- +goose StatementEnd
