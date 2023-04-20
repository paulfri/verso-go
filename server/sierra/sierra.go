package sierra

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/citra/server/utils"
)

type SierraRouter struct {
	Controller utils.Controller
}

func Router(controller utils.Controller) http.Handler {
	sierra := SierraRouter{Controller: controller}

	router := chi.NewRouter()
	router.Post("/api/0/accounts/ClientLogin", sierra.token)
	router.Get("/api/0/status", sierra.status)
	router.Get("/api/0/user-info", sierra.user)

	return router
}
