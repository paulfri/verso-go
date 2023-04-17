-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id serial NOT NULL,
    uuid uuid NOT NULL DEFAULT gen_random_uuid(),
    title text NOT NULL,
    url text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),

    PRIMARY KEY(id),
    UNIQUE(uuid)
);

CREATE TRIGGER feeds_touch_updated_at
BEFORE UPDATE ON feeds
  FOR EACH ROW EXECUTE PROCEDURE touch_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
