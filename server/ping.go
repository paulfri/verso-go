package server

import (
	"errors"
	"net/http"

	"github.com/unrolled/render"
)

type PingResponse struct {
	Status string `json:"status"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	render.New().JSON(w, http.StatusOK, PingResponse{Status: "ok"})
}

func testPanic(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("test panic"))
}
