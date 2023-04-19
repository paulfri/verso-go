package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/versolabs/citra/db/query"
)

func (c *Controller) feedIndex(ctx *gin.Context) {
	feeds, err := c.queries.ListRSSFeeds(ctx)

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

func (c *Controller) feedShow(ctx *gin.Context) {
	uuid := uuid.Must(uuid.Parse(ctx.Param("pk")))
	feed, err := c.queries.GetRssFeedByUuid(ctx, uuid)

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

type CreateRssFeedRequest struct {
	Title string `json:"title" binding:"required"`
	Url   string `json:"url" binding:"required,http_url"`
}

func (c *Controller) feedCreate(ctx *gin.Context) {
	var req CreateRssFeedRequest

	if err := ctx.BindJSON(&req); err != nil {
		fmt.Println(err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})

		return
	}

	feed, err := c.queries.CreateRssFeed(
		ctx,
		query.CreateRssFeedParams{Title: req.Title, Url: req.Url},
	)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"feed": feed,
		})
	} else {
		fmt.Println(err)

		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err,
			"feed":  nil,
		})
	}
}
