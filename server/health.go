package server

import (
	"net/http"

	"github.com/unrolled/render"
)

func ping(w http.ResponseWriter, r *http.Request) {
	render.New().JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
