-- +goose Up
create extension if not exists pgcrypto;

-- +goose Down
drop extension pgcrypto;
