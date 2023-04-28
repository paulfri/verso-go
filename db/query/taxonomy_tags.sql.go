// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: taxonomy_tags.sql

package query

import (
	"context"
)

const createTag = `-- name: CreateTag :one
insert into taxonomy.tags (
  name,
  user_id
)
values ($1, $2)
returning id, uuid, created_at, updated_at, user_id, name
`

type CreateTagParams struct {
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

func (q *Queries) CreateTag(ctx context.Context, arg CreateTagParams) (TaxonomyTag, error) {
	row := q.db.QueryRowContext(ctx, createTag, arg.Name, arg.UserID)
	var i TaxonomyTag
	err := row.Scan(
		&i.ID,
		&i.UUID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
	)
	return i, err
}

const deleteTagByNameAndUserID = `-- name: DeleteTagByNameAndUserID :exec
delete from taxonomy.tags
where
  name = $1
  and user_id = $2
`

type DeleteTagByNameAndUserIDParams struct {
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
}

func (q *Queries) DeleteTagByNameAndUserID(ctx context.Context, arg DeleteTagByNameAndUserIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteTagByNameAndUserID, arg.Name, arg.UserID)
	return err
}

const getTagsByUserID = `-- name: GetTagsByUserID :many
select id, uuid, created_at, updated_at, user_id, name from taxonomy.tags
where user_id = $1
`

func (q *Queries) GetTagsByUserID(ctx context.Context, userID int64) ([]TaxonomyTag, error) {
	rows, err := q.db.QueryContext(ctx, getTagsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TaxonomyTag
	for rows.Next() {
		var i TaxonomyTag
		if err := rows.Scan(
			&i.ID,
			&i.UUID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const renameTagByNameAndUserID = `-- name: RenameTagByNameAndUserID :exec
update taxonomy.tags
  set name = $3::text
where
  name = $1
  and user_id = $2
`

type RenameTagByNameAndUserIDParams struct {
	Name    string `json:"name"`
	UserID  int64  `json:"user_id"`
	NewName string `json:"new_name"`
}

func (q *Queries) RenameTagByNameAndUserID(ctx context.Context, arg RenameTagByNameAndUserIDParams) error {
	_, err := q.db.ExecContext(ctx, renameTagByNameAndUserID, arg.Name, arg.UserID, arg.NewName)
	return err
}
