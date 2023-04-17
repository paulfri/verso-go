package http

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
)

type Controller struct {
	queries *query.Queries
}

func Serve() {
	r := gin.Default()

	controller := Controller{
		queries: db.Queries(),
	}

	r.GET("/ping", controller.ping)
	r.GET("/feeds", controller.feedIndex)
	r.GET("/feeds/:pk", controller.feedShow)
	r.POST("/feeds", controller.feedCreate)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	r.Run(bind)
}
