-- name: GetUser :one
select * from identity.users where id = $1;

-- name: GetUserByEmail :one
select * from identity.users where email = $1;

-- name: CreateUser :one
insert into identity.users (
  email,
  name,
  password,
  superuser
) values ($1, $2, $3, $4) returning *;
