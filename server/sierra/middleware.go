package sierra

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ContextKey string

const ContextUserIDKey ContextKey = "user_id"

func (r SierraRouter) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		split := strings.Split(auth, "GoogleLogin auth=")

		if len(split) != 2 {
			authFailed(w)
			return
		} else {
			identifier := split[1]

			if identifier != "" {
				fmt.Println(identifier)

				token, err := r.Controller.Queries.GetTokenByIdentifier(req.Context(), identifier)

				if err != nil {
					authFailed(w)
					return
				}

				ctx := context.WithValue(req.Context(), ContextUserIDKey, token.UserID)

				next.ServeHTTP(w, req.WithContext(ctx))
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
