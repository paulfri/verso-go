-- +goose NO TRANSACTION
-- +goose Up
create table taxonomy.rss_feed_tags (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  tag_id bigint not null references taxonomy.tags(id) on delete cascade,
  rss_feed_id bigint not null references rss.feeds(id) on delete cascade
);

create unique index concurrently if not exists taxonomy_rss_feed_tags_tag_id_rss_feed_id_key
  on taxonomy.rss_feed_tags (tag_id, rss_feed_id);

create trigger taxonomy_rss_feed_tags_touch_updated_at
  before update on taxonomy.rss_feed_tags for each row
  execute procedure touch_updated_at();

-- +goose Down
drop table taxonomy.rss_feed_tags;
