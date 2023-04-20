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
		auth.Get("/api/0/status", rainier.MetaStatus)

		// auth & identity
		auth.Get("/api/0/ping", rainier.MetaPing)
		auth.Get("/api/0/token", rainier.UserTokenGet)
		auth.Get("/api/0/user-info", rainier.UserGet)

		// subscriptions
		auth.Post("/api/0/subscription/quickadd", rainier.SubscriptionCreate)
	})

	return router
}
