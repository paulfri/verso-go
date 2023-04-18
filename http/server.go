package http

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
)

type Controller struct {
	asynq   *asynq.Client
	queries *query.Queries
}

const redisAddr = "127.0.0.1:6379"

func Serve(cliContext *cli.Context) error {
	r := gin.Default()

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	controller := Controller{
		queries: db.Queries(),
		asynq:   asynqClient,
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
