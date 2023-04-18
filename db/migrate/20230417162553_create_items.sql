-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id serial NOT NULL,
  uuid uuid NOT NULL DEFAULT gen_random_uuid(),
  feed_id int NOT NULL REFERENCES feeds(id),
  remote_id text NOT NULL,
  title text NOT NULL,
  link text NOT NULL,
  content text NOT NULL,
  published_at timestamp with time zone,
  remote_updated_at timestamp with time zone,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now(),

  PRIMARY KEY(id),
  UNIQUE(uuid)
);

CREATE TRIGGER items_touch_updated_at
BEFORE UPDATE ON items
  FOR EACH ROW EXECUTE PROCEDURE touch_updated_at();

CREATE UNIQUE INDEX items_remote_id_feed_id_index
  ON items (feed_id, remote_id);

CREATE INDEX items_published_at_index
  ON items (published_at DESC);

CREATE INDEX items_updated_at_index
  ON items (updated_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
