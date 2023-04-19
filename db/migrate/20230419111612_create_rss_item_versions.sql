-- +goose Up
-- +goose StatementBegin
create table content.rss_item_versions (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  rss_item_id bigint references content.rss_items(id) on delete cascade,

  title text not null,
  link text not null,
  content text not null,
  published_at timestamp with time zone,
  remote_updated_at timestamp with time zone
);

create index content_rss_item_versions_rss_item_id_fkey on content.rss_item_versions (rss_item_id);

create trigger content_rss_item_versions_touch_updated_at
  before update on content.rss_item_versions for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table content.rss_item_versions;
-- +goose StatementEnd
