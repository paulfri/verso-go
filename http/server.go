package http

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/http/sierra"
	"github.com/versolabs/citra/tasks"
)

type Controller struct {
	asynq   *asynq.Client
	queries *query.Queries
}

func Serve(cliContext *cli.Context) error {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)

	controller := Controller{
		queries: db.Queries(),
		asynq:   tasks.Client(),
	}

	engine.GET("/ping", controller.ping)
	engine.GET("/feeds", controller.feedIndex)
	engine.GET("/feeds/:pk", controller.feedShow)
	engine.POST("/feeds", controller.feedCreate)
	// TODO: remove
	engine.POST("/crawl", controller.crawl)

	sierra.Routes(engine)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	engine.Run(bind)

	return nil
}
