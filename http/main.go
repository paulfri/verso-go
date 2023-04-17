package http

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

type Controller struct {
	queries *query.Queries
}

func Serve() {
	r := gin.Default()

	// TODO: parameterize
	database, err := sql.Open("postgres", "user=citra dbname=citra_dev sslmode=disable")
	queries := query.New(database)
	if err != nil {
		fmt.Println(err)
	}

	controller := Controller{
		queries: queries,
	}

	r.GET("/ping", controller.ping)
	r.GET("/feeds", controller.feedIndex)
	r.GET("/feeds/:pk", controller.feedShow)
	r.POST("/feeds", controller.feedCreate)

	// TODO: parameterize
	r.Run("localhost:8080")
}
