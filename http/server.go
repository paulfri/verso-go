package http

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/tasks"
)

type Controller struct {
	asynq   *asynq.Client
	queries *query.Queries
}

func Serve(cliContext *cli.Context) error {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	controller := Controller{
		queries: db.Queries(),
		asynq:   tasks.Client(),
	}

	r.GET("/ping", controller.ping)
	r.GET("/feeds", controller.feedIndex)
	r.GET("/feeds/:pk", controller.feedShow)
	r.POST("/feeds", controller.feedCreate)

	// TODO: remove
	r.POST("/crawl", controller.crawl)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	r.Run(bind)

	return nil
}
