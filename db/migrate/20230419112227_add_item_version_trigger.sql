-- +goose Up
-- +goose StatementBegin
create function create_rss_item_version()
returns trigger as $$
begin
  insert into content.rss_item_versions (
    rss_item_id,
    title,
    link,
    content,
    published_at,
    remote_updated_at
  ) select old.id, old.title, old.link, old.content, old.published_at, old.remote_updated_at
    where exists (
      select 1 from content.rss_items i where
        i.id = old.id and (
          i.title <> new.title or
          i.link <> new.link or
          i.content <> new.content or
          i.published_at <> new.published_at or
          i.remote_updated_at <> new.remote_updated_at
        )
    );

  return new;
end;
$$ language 'plpgsql';

create trigger rss_item_versions_insert_on_item_update
  before update on content.rss_items for each row
  execute procedure create_rss_item_version();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger rss_item_versions_insert_on_item_update on content.rss_items;
drop function create_rss_item_version();
-- +goose StatementEnd
