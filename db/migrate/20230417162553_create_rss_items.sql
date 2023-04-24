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
  reader_id bigint unique not null
);

create trigger rss_items_touch_updated_at
  before update on rss.items for each row
  execute procedure touch_updated_at();

create unique index rss_items_rss_feed_id_rss_guid_key
  on rss.items (feed_id, rss_guid);
create index rss_items_published_at_index
  on rss.items (published_at desc);

-- Reader compatibility requires an int64 identifier. This insert trigger
-- generates a new one derived from the UUID rather than expose the serial ID.
create function generate_reader_id()
returns trigger as $$
begin
  new.reader_id = ('x' || translate(new.uuid::text, '-', ''))::bit(64)::bigint;

  return new;
end;
$$ language plpgsql;

create trigger rss_items_generate_reader_id
  before insert on rss.items for each row
  execute procedure generate_reader_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete function generate_reader_id();
drop table rss.items;
-- +goose StatementEnd
