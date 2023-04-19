package sierra

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c SierraController) user(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"userId":              "00157a17b192950b65be3791",
		"userName":            "Paul Friedman",
		"userProfileId":       "00157a17b192950b65be3791",
		"userEmail":           "paul@verso.so",
		"isBloggerUser":       false,
		"signupTimeSec":       1370709105,
		"isMultiLoginEnabled": false,
		"isPremium":           true,
	})
}
