package utils

import (
	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/versolabs/citra/db/query"
)

type Container struct {
	Asynq   *asynq.Client
	Queries *query.Queries
	Render  *render.Render
}