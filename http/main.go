package http

import (
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

	// TODO: parameterize
	r.Run("localhost:8080")
}
