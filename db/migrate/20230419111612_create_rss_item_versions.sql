-- +goose NO TRANSACTION
-- +goose Up
create table rss.item_versions (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  item_id bigint references rss.items(id) on delete cascade,

  title text not null,
  link text not null,
  content text not null,
  published_at timestamptz,
  remote_updated_at timestamptz
);

create index concurrently if not exists rss_item_versions_item_id_fkey
  on rss.item_versions (item_id);

create trigger rss_item_versions_touch_updated_at
  before update on rss.item_versions for each row
  execute procedure touch_updated_at();

-- +goose Down
drop table rss.item_versions;
