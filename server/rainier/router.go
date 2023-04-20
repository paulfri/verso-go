package rainier

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/citra/server/utils"
)

type RainierController struct {
	Container *utils.Container
}

func Router(container *utils.Container) http.Handler {
	rainier := RainierController{Container: container}

	router := chi.NewRouter()
	router.Post("/api/0/accounts/ClientLogin", rainier.login)

	router.With(rainier.AuthMiddleware).Route("/", func(auth chi.Router) {
		auth.Get("/api/0/status", rainier.status)

		// auth & identity
		auth.Get("/api/0/token", rainier.token)
		auth.Get("/api/0/user-info", rainier.user)
	})

	return router
}
