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
	"github.com/versolabs/verso/config"
	"github.com/versolabs/verso/core/command"
	"github.com/versolabs/verso/db"
	vm "github.com/versolabs/verso/middleware"
	"github.com/versolabs/verso/server/reader"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker"
)

func Serve(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		router := Router(config)
		bind := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
		http.ListenAndServe(bind, router)

		return nil
	}
}

func Router(config *config.Config) *chi.Mux {
	airbrake := util.Airbrake(config)
	logger := util.Logger()
	asynq := worker.Client(config.RedisURL)
	db, queries := db.Init(config.Database.URL(), config.Database.Migrate)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(vm.LoggerMiddleware(logger))
	r.Use(middleware.Timeout(60 * time.Second))

	if config.Debug {
		r.Use(vm.DebugResponseBody(logger))
		r.Use(vm.DebugRequestBody(logger))
	}

	// Error notifier middleware goes last.
	r.Use(vm.NotifyAirbrake(airbrake))

	command := &command.Command{
		Asynq:   asynq,
		DB:      db,
		Queries: queries,
	}

	container := util.Container{
		Airbrake:  airbrake,
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

	r.Mount("/", reader.LoginRouter(&container))
	r.Mount("/reader/api/0", reader.Router(&container))

	return r
}
