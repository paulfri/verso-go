package rainier

import (
	"net/http"
)

type StatusValue string

const (
	STATUS_OK  StatusValue = "ok"
	STATUS_ERR StatusValue = "down"
)

const REDIRECT_URL = "https://verso.so"
const ERROR_TEXT = "An unexpected error occurred."

type StatusResponse struct {
	Status      StatusValue `json:"status"`
	Description string      `json:"description,omitempty"`
	Redirect    string      `json:"redirect,omitempty"`
}

type None struct{}

func (c *RainierController) MetaStatus(w http.ResponseWriter, req *http.Request) {
	var err *None

	if err != nil {
		c.Container.Render.JSON(w, http.StatusOK, StatusResponse{
			Status:      STATUS_ERR,
			Description: ERROR_TEXT,
			Redirect:    REDIRECT_URL,
		})
	} else {
		c.Container.Render.JSON(w, http.StatusOK, StatusResponse{
			Status: STATUS_OK,
		})
	}
}

func (c *RainierController) MetaPing(w http.ResponseWriter, req *http.Request) {
	c.Container.Render.Text(w, http.StatusOK, "OK")
}
