package sierra

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c SierraController) status(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}
