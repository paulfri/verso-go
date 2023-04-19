-- +goose Up
-- +goose StatementBegin
create extension if not exists pgcrypto;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop extension pgcrypto;
-- +goose StatementEnd
