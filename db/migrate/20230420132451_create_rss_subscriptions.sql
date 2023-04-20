-- +goose Up
-- +goose StatementBegin
create table content.rss_subscriptions (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint references identity.users(id) not null,
  rss_feed_id bigint references content.rss_feeds(id) not null,
  custom_title text
);

create trigger content_rss_subscriptions_touch_updated_at
  before update on content.rss_subscriptions for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table content.rss_subscriptions;
-- +goose StatementEnd
