package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// Debug middleware to log the HTTP request body.
func DebugRequestBody(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error().Msg(err.Error())
			}

			bodyStr := string(body)
			if bodyStr != "" {
				logger.Debug().Str("request_body", bodyStr).Msg("")
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			next.ServeHTTP(w, r)
		})
	}
}

type responseBodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Debug middleware to log the HTTP response body. This is accomplished by
// swapping out the default http.ResponseWriter with a struct that writes to
// both the default, and an additional buffer.
func DebugResponseBody(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body := &bytes.Buffer{}
			newW := &responseBodyWriter{body: body, ResponseWriter: w}
			next.ServeHTTP(newW, r)

			bodyRead, err := io.ReadAll(body)
			if err != nil {
				logger.Error().Msg(err.Error())
			}

			bodyStr := string(bodyRead)
			if bodyStr != "" {
				logger.Debug().Str("response_body", bodyStr).Msg("")
			}
		})
	}
}

// Middleware to log request metadata.
// Adpated from https://github.com/ironstar-io/chizerolog
func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				t2 := time.Now()

				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// TODO: this isn't working
				id := middleware.GetReqID(r.Context())

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
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
