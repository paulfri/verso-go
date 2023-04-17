package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Controller) feedShow(ctx *gin.Context) {
	uuid := uuid.Must(uuid.Parse(ctx.Param("pk")))
	feed, err := c.queries.GetFeed(ctx, uuid)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"error": nil,
			"feed":  feed,
		})
	} else {
		fmt.Println(err)

		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
			"feed":  nil,
		})
	}
}

func (c *Controller) feedIndex(ctx *gin.Context) {
	feeds, err := c.queries.ListFeeds(ctx)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"feeds": feeds,
		})
	} else {
		fmt.Println(err)

		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err,
			"feed":  nil,
		})
	}
}
