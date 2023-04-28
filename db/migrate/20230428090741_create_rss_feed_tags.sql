-- +goose Up
-- +goose StatementBegin
create table taxonomy.rss_feed_tags (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  tag_id bigint not null references taxonomy.tags(id) on delete cascade,
  rss_feed_id bigint not null references rss.feeds(id) on delete cascade
);

create trigger taxonomy_rss_feed_tags_touch_updated_at
  before update on taxonomy.rss_feed_tags for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table taxonomy.rss_feed_tags;
-- +goose StatementEnd
