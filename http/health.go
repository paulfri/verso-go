package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
