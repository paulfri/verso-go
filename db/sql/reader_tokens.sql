-- name: CreateReaderToken :one
insert into identity.reader_tokens (user_id, identifier)
  values ($1, $2)
  returning *;

-- name: GetReaderTokenByIdentifier :one
select * from identity.reader_tokens where identifier = $1;
