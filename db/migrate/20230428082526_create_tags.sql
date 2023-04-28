-- +goose NO TRANSACTION
-- +goose Up
create schema taxonomy;
create table taxonomy.tags (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  user_id bigint not null references identity.users(id) on delete cascade,

  name text not null
);

create unique index concurrently if not exists taxonomy_tags_user_id_name_key
  on taxonomy.tags (user_id, name);

create trigger taxonomy_tags_touch_updated_at
  before update on taxonomy.tags for each row
  execute procedure touch_updated_at();

-- +goose Down
drop table taxonomy.tags;
drop schema taxonomy;
