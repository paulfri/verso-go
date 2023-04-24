package server

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func PrintRequestBody(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(body))

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}

func PrintRequestBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PrintRequestBody(w, r)
		next.ServeHTTP(w, r)
	})
}
