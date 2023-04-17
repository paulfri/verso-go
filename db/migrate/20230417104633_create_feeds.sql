-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id    int NOT NULL,
    title text,
    url text,
    PRIMARY KEY(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
