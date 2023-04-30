package middleware

import (
	"net/http"
	"time"

	cm "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// Middleware to log request metadata.
// Adpated from https://github.com/ironstar-io/chizerolog
func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()
			ww := cm.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				t2 := time.Now()

				// TODO: this isn't working
				id := cm.GetReqID(r.Context())

				log.Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"id":         id,
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
						"query":      r.URL.Query(),
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
