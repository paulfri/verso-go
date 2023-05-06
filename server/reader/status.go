package reader

import (
	"net/http"
)

type StatusValue string

const (
	StatusOk  StatusValue = "ok"
	StatusErr StatusValue = "down"
)

const RedirectURL = "https://verso.so"
const ErrorText = "An unexpected error occurred."

type StatusResponse struct {
	Status      StatusValue `json:"status"`
	Description string      `json:"description,omitempty"`
	Redirect    string      `json:"redirect,omitempty"`
}

type None struct{}

func (c *Controller) MetaStatus(w http.ResponseWriter, _ *http.Request) {
	var err *None

	if err != nil {
		c.Container.Render.JSON(w, http.StatusOK, StatusResponse{
			Status:      StatusErr,
			Description: ErrorText,
			Redirect:    RedirectURL,
		})
	} else {
		c.Container.Render.JSON(w, http.StatusOK, StatusResponse{
			Status: StatusOk,
		})
	}
}

func (c *Controller) MetaPing(w http.ResponseWriter, _ *http.Request) {
	c.Container.Render.Text(w, http.StatusOK, "OK")
}
