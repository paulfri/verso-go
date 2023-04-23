package reader

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/verso/util"
)

type ReaderController struct {
	Container *util.Container
}

// This handler allows the main server to mount the ClientLogin endpoint
// at the API root instead of under the /reader/api/0 path.
func LoginRouter(container *util.Container) http.Handler {
	reader := ReaderController{Container: container}

	router := chi.NewRouter()
	router.Post("/accounts/ClientLogin", reader.ClientLogin)

	return router
}

func Router(container *util.Container) http.Handler {
	reader := ReaderController{Container: container}

	router := chi.NewRouter()
	router.Get("/status", reader.MetaStatus)
	router.Post("/accounts/ClientLogin", reader.ClientLogin)

	router.With(reader.AuthMiddleware).Route("/", func(auth chi.Router) {
		auth.Get("/ping", reader.MetaPing)

		// auth & identity
		auth.Get("/token", reader.UserToken)
		auth.Get("/user-info", reader.User)
		auth.Get("/preference/list", reader.UserPreferences)
		auth.Get("/preference/stream/list", reader.UserStreamPreferences)
		auth.Get("/user/friend/list", reader.UserFriendList)

		// subscriptions
		auth.Post("/subscription/quickadd", reader.SubscriptionQuickAdd)
		auth.Get("/subscribed", reader.SubscriptionExists)

		// stream
		auth.Get("/stream/contents/*", reader.StreamContents)
	})

	return router
}
