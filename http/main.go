package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/versolabs/citra/db/query"
)

func Serve() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/feeds", func(c *gin.Context) {
		ctx := context.Background()

		database, err := sql.Open("postgres", "user=citra dbname=citra_dev sslmode=disable")
		queries := query.New(database)
		if err != nil {
			fmt.Println(err)
		}

		// list all authors
		feeds, err := queries.ListFeeds(ctx)
		fmt.Println(feeds)
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"feeds": feeds,
		})
	})

	r.Run()
}
