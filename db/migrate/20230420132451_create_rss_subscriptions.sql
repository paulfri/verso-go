-- +goose Up
-- +goose StatementBegin
create table rss.subscriptions (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint not null references identity.users(id) on delete cascade,
  feed_id bigint not null references rss.feeds(id) on delete cascade,
  custom_title text
);

create index rss_subscriptions_user_id_fkey
  on rss.subscriptions (user_id);

create index rss_subscriptions_feed_id_fkey
  on rss.subscriptions (feed_id);

create unique index rss_subscriptions_user_id_feed_id_key
  on rss.subscriptions (user_id, feed_id);

create trigger rss_subscriptions_touch_updated_at
  before update on rss.subscriptions for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rss.subscriptions;
-- +goose StatementEnd
