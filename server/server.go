package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/server/rainier"
	"github.com/versolabs/citra/tasks"
	"github.com/versolabs/citra/util"
)

func Serve(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(60 * time.Second))

		db, queries := db.Init(config.DatabaseUrl)

		container := util.Container{
			Asynq:   tasks.Client(config.RedisUrl),
			DB:      db,
			Queries: queries,
			Render:  render.New(),
		}

		r.Get("/ping", ping)
		r.Mount("/reader/api/0", rainier.Router(&container))

		bind := fmt.Sprintf("%s:%s", config.Host, config.Port)
		http.ListenAndServe(bind, r)

		return nil
	}
}
