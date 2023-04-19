package sierra

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/tasks"
)

// TODO: DRY with http.Controller
type SierraController struct {
	asynq   *asynq.Client
	queries *query.Queries
}

func Routes(engine *gin.Engine) {
	sierra := engine.Group("/reader")

	controller := SierraController{
		queries: db.Queries(),
		asynq:   tasks.Client(),
	}

	sierra.GET("/api/0/status", controller.status)
	sierra.GET("/api/0/user-info", controller.user)
}
