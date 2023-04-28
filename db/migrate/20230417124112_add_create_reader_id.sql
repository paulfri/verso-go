-- +goose Up
-- +goose StatementBegin
create or replace function generate_reader_id()
returns trigger as $$
begin
  -- reader_id is stored as a 16-length zero-added hex representation of a
  -- 64-bit integer (bigint)
  new.reader_id = lpad(
    -- convert to hex
    to_hex(
      -- floor to convert the random numeric to whole number
      floor(
        -- 9223372036854775807 is the maximum value of a signed 64-bit integer
        -- random() is cast to numeric for greater precision (otherwise it
        -- reduces the number of possible random bigints)
        random()::numeric * 9223372036854775807
      -- cast to bigint for insertion
      )::bigint),
    16,
    '0'
  );

  return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down
drop function generate_reader_id();
