package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/core/command"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/server/reader"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker"
)

func Serve(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		logger := util.Logger()

		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		r.Use(LoggerMiddleware(logger))
		r.Use(middleware.Timeout(60 * time.Second))

		if config.Debug {
			r.Use(DebugResponseBody(logger))
			r.Use(DebugRequestBody(logger))
		}

		asynq := worker.Client(config.RedisURL)
		db, queries := db.Init(config.DatabaseURL, config.DatabaseMigrate)

		command := &command.Command{
			Asynq:   asynq,
			DB:      db,
			Queries: queries,
		}

		container := util.Container{
			Asynq:     asynq,
			Command:   command,
			Config:    config,
			DB:        db,
			Logger:    logger,
			Queries:   queries,
			Render:    render.New(),
			Validator: validator.New(),
		}

		r.Get("/ping", ping)
		r.Get("/panic", testPanic)

		r.Mount("/", reader.LoginRouter(&container))
		r.Mount("/reader/api/0", reader.Router(&container))

		bind := fmt.Sprintf("%s:%s", config.Host, config.Port)
		http.ListenAndServe(bind, r)

		return nil
	}
}
