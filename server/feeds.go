package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/versolabs/citra/server/utils"
)

type FeedsRouter struct {
	controller utils.Controller
}

func NewFeedsRouter(controller utils.Controller) http.Handler {
	feeds := FeedsRouter{controller: controller}
	r := chi.NewRouter()

	r.Get("/", feeds.feedIndex)
	r.Get("/{id}", feeds.feedShow)

	return r
}

func (c *FeedsRouter) feedIndex(w http.ResponseWriter, req *http.Request) {
	feeds, err := c.controller.Queries.ListRSSFeeds(req.Context())

	if err == nil {
		c.controller.Render.JSON(w, http.StatusOK, map[string]any{"feeds": feeds})
	} else {
		fmt.Println(err)

		c.controller.Render.JSON(w, http.StatusInternalServerError, map[string]any{
			"error": err,
			"feed":  nil,
		})
	}
}

func (c *FeedsRouter) feedShow(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	val := chi.URLParam(req, "id")
	uuid := uuid.Must(uuid.Parse(val))
	feed, err := c.controller.Queries.GetRssFeedByUuid(ctx, uuid)

	if err == nil {
		c.controller.Render.JSON(w, http.StatusOK, map[string]any{
			"error": nil,
			"feed":  feed,
		})
	} else {
		fmt.Println(err)

		c.controller.Render.JSON(w, http.StatusInternalServerError, map[string]any{
			"error": "not_found", // TODO: constantize
			"feed":  nil,
		})
	}
}
