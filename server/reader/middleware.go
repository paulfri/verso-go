package reader

import (
	"context"
	"net/http"
	"strings"
)

type ContextUserIDKey struct{}
type ContextAuthTokenKey struct{}

func (c *ReaderController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		auth := req.Header.Get("Authorization")
		split := strings.Split(auth, "GoogleLogin auth=")

		if len(split) != 2 {
			authFailed(w)
			return
		} else {
			identifier := split[1]

			if identifier != "" {
				queries := c.Container.GetQueries(req)
				token, err := queries.GetReaderTokenByIdentifier(
					ctx,
					identifier,
				)

				if err != nil {
					authFailed(w)
					return
				}

				withUser := context.WithValue(ctx, ContextUserIDKey{}, token.UserID)
				withToken := context.WithValue(withUser, ContextAuthTokenKey{}, token.Identifier)

				next.ServeHTTP(w, req.WithContext(withToken))
			} else {
				authFailed(w)
				return
			}
		}
	})
}

func authFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte("Unauthorized"))
}
