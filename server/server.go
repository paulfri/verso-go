package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/server/sierra"
	"github.com/versolabs/citra/server/utils"
	"github.com/versolabs/citra/tasks"
)

func Serve(cliContext *cli.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	controller := utils.Controller{
		Queries: db.Queries(),
		Asynq:   tasks.Client(),
		Render:  render.New(),
	}

	r.Get("/ping", ping)
	r.Mount("/feeds", NewFeedsRouter(controller))
	r.Mount("/reader", sierra.Router(controller))

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	http.ListenAndServe(bind, r)

	return nil
}
