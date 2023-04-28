-- +goose Up
-- +goose StatementBegin
create table identity.reader_tokens (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint references identity.users(id) not null,
  identifier text unique not null,
  revoked_at timestamp
);

create index identity_reader_tokens_user_id_fkey
  on identity.reader_tokens (user_id);

create trigger identity_reader_tokens_touch_updated_at
  before update on identity.reader_tokens for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table identity.reader_tokens;
-- +goose StatementEnd
