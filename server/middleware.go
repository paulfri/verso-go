package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
)

// Debug middleware to log the HTTP request body.
func DebugRequestBodyMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
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
func DebugResponseBodyMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
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
