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

-- name: CreateToken :one
insert into identity.tokens (user_id, identifier) values ($1, $2) returning *;

-- name: GetTokenByIdentifier :one
select * from identity.tokens where identifier = $1;
