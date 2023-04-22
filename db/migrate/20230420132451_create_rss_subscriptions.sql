-- +goose Up
-- +goose StatementBegin
create table rss.subscriptions (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint references identity.users(id) not null,
  feed_id bigint references rss.feeds(id) not null,
  custom_title text
);

create trigger rss_subscriptions_touch_updated_at
  before update on rss.subscriptions for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rss.subscriptions;
-- +goose StatementEnd
