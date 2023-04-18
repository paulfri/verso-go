-- +goose Up
-- +goose StatementBegin
CREATE TABLE items (
  id serial NOT NULL,
  uuid uuid NOT NULL DEFAULT gen_random_uuid(),
  feed_id int NOT NULL REFERENCES feeds(id),
  title text,
  link text,
  content text,
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

CREATE INDEX items_published_at_index
  ON items (published_at DESC);

CREATE INDEX items_updated_at_index
  ON items (updated_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
-- +goose StatementEnd
