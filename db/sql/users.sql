-- name: CreateUser :one
insert into identity.users (
  email,
  name,
  password,
  superuser
) values ($1, $2, $3, $4) returning *;

-- name: GetUserByEmail :one
select * from identity.users where email = $1;

-- name: GetUserById :one
select * from identity.users where id = $1;

-- name: CreateReaderToken :one
insert into identity.reader_tokens (user_id, identifier) values ($1, $2) returning *;

-- name: GetReaderTokenByIdentifier :one
select * from identity.reader_tokens where identifier = $1;
