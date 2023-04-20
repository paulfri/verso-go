-- +goose Up
-- +goose StatementBegin
create schema identity;

create table identity.users (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  email text not null unique,
  name text not null,
  password text,
  superuser bool not null default false
);

create trigger identity_users_touch_updated_at
  before update on identity.users for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table identity.users;
drop schema identity;
-- +goose StatementEnd
