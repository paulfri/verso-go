-- +goose Up
-- +goose StatementBegin
create table rss.items (
  id bigint primary key generated always as identity,
  uuid uuid unique not null default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),

  feed_id bigint not null references rss.feeds(id) on delete cascade,

  rss_guid text not null,
  title text not null,
  link text not null,
  author text,
  author_email text,
  content text not null,
  summary text,
  published_at timestamp with time zone,
  remote_updated_at timestamp with time zone,
  reader_id text unique not null constraint reader_id_length_check
    check (char_length(reader_id) = 16)
);

create unique index rss_items_rss_feed_id_rss_guid_key
  on rss.items (feed_id, rss_guid);

create index rss_items_published_at_idx
  on rss.items (published_at desc);

create function generate_reader_id()
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

create trigger rss_items_generate_reader_id
  before insert on rss.items for each row
  execute procedure generate_reader_id();

create trigger rss_items_touch_updated_at
  before update on rss.items for each row
  execute procedure touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete function generate_reader_id();
drop table rss.items;
-- +goose StatementEnd
