package server

import (
	"net/http"
	"os"

	"github.com/unrolled/render"
)

type PingResponse struct {
	Status   string `json:"status"`
	Revision string `json:"revision"`
}

func ping(w http.ResponseWriter, _ *http.Request) {
	revision := os.Getenv("RENDER_GIT_COMMIT")
	if revision != "" {
		revision = "unknown"
	}

	render.New().JSON(w, http.StatusOK, PingResponse{
		Status:   "ok",
		Revision: revision,
	})
}
