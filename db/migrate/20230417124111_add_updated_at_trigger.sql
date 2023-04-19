-- +goose Up
-- +goose StatementBegin
create function touch_updated_at()   
returns trigger as $$
begin
  new.updated_at = now();
  return new;   
end;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function touch_updated_at();
-- +goose StatementEnd
