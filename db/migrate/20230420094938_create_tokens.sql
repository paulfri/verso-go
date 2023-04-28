-- +goose NO TRANSACTION
-- +goose Up
create table identity.reader_tokens (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  user_id bigint not null references identity.users(id) on delete cascade,
  identifier text unique not null,
  revoked_at timestamptz
);

create index concurrently if not exists identity_reader_tokens_user_id_fkey
  on identity.reader_tokens (user_id);

create trigger identity_reader_tokens_touch_updated_at
  before update on identity.reader_tokens for each row
  execute procedure touch_updated_at();

-- +goose Down
drop table identity.reader_tokens;
