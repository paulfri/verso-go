package util

import (
	"database/sql"

	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/versolabs/citra/core/command"
	"github.com/versolabs/citra/db/query"
)

type Container struct {
	Asynq   *asynq.Client
	Command *command.Command
	DB      *sql.DB
	Queries *query.Queries
	Render  *render.Render
}
