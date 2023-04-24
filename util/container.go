package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hetiansu5/urlquery"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
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
	Logger    *zerolog.Logger
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

func (c Container) Form(req *http.Request, s interface{}) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}

	err = urlquery.Unmarshal([]byte(req.Form.Encode()), s)
	if err != nil {
		return err
	}

	return c.Validator.Struct(s)
}

func (c Container) JSONBody(req *http.Request, s interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(s)

	if err != nil {
		return err
	}

	return nil
}

// Given a struct with request parameters, unmarshal the query string from the
// given request into that struct.
func (c Container) BodyParams(s interface{}, req *http.Request) error {
	body, _ := ioutil.ReadAll(req.Body)
	asByte := []byte(body)

	fmt.Printf("%v\n", body)

	err := urlquery.Unmarshal(asByte, s)
	if err != nil {
		return err
	}

	return c.Validator.Struct(s)
}
