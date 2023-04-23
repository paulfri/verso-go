package util

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hetiansu5/urlquery"
	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/versolabs/verso/core/command"
	"github.com/versolabs/verso/db/query"
)

type Container struct {
	Asynq     *asynq.Client
	Command   *command.Command
	Config    *Config
	DB        *sql.DB
	Queries   *query.Queries
	Render    *render.Render
	Validator *validator.Validate
}

// Given a struct with request parameters, unmarshal the query string from the
// given request into that struct.
func (c Container) Params(s interface{}, req *http.Request) error {
	query := req.URL.RawQuery
	asByte := []byte(query)

	err := urlquery.Unmarshal(asByte, s)
	if err != nil {
		return err
	}

	return c.Validator.Struct(s)
}
