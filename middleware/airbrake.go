package middleware

import (
	"net/http"

	"github.com/airbrake/gobrake/v5"
)

func NotifyAirbrake(airbrake *gobrake.Notifier) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					// Notify airbrake.
					airbrake.Notify(rvr, r)

					// Reraise to be handled by recoverer middleware.
					panic(rvr)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
