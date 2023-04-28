-- name: CreateTag :one
insert into taxonomy.tags (
  name,
  user_id
)
values ($1, $2)
returning *;

-- name: GetTagsByUserID :many
select * from taxonomy.tags
where user_id = $1;

-- name: RenameTagByNameAndUserID :exec
update taxonomy.tags
  set name = @new_name::text
where
  name = $1
  and user_id = $2;

-- name: DeleteTagByNameAndUserID :exec
delete from taxonomy.tags
where
  name = $1
  and user_id = $2;
