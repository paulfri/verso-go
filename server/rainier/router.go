package rainier

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/citra/util"
)

type RainierController struct {
	Container *util.Container
}

func Router(container *util.Container) http.Handler {
	rainier := RainierController{Container: container}

	router := chi.NewRouter()
	router.Get("/status", rainier.MetaStatus)
	router.Post("/accounts/ClientLogin", rainier.login)

	router.With(rainier.AuthMiddleware).Route("/", func(auth chi.Router) {
		auth.Get("/ping", rainier.MetaPing)

		// auth & identity
		auth.Get("/token", rainier.UserToken)
		auth.Get("/user-info", rainier.User)
		auth.Get("/preference/list", rainier.UserPreferences)
		auth.Get("/preference/stream/list", rainier.UserStreamPreferences)
		auth.Get("/user/friend/list", rainier.UserFriendList)

		// subscriptions
		auth.Post("/subscription/quickadd", rainier.SubscriptionQuickAdd)

		// stream
		auth.Get("/stream/contents/*", rainier.StreamContents)
	})

	return router
}
