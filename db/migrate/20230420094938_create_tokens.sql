-- +goose Up
-- +goose StatementBegin
create table identity.tokens (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  user_id bigint references identity.users(id) not null,
  identifier text unique not null,
  revoked_at timestamp
);

create trigger identity_tokens_touch_updated_at
  before update on identity.tokens for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table identity.tokens;
-- +goose StatementEnd
