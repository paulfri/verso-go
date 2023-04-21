package rainier

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SortOrderValue string

const (
	Ascending  SortOrderValue = "o"
	Descending SortOrderValue = "a"
)

type StreamContentsRequestParams struct {
	// r: sort criteria. Items are sorted by date (descending by default), r=o inverts the order.
	SortOrder SortOrderValue `query:"r"`
	// n: the number of items per page. Default: 20.
	NumItems int `query:"n"`
	// c: the continuation string (see below).
	Continuation string `query:"c"`
	// xt: a stream ID to exclude from the list.
	Exclude string `query:"xt"`
	// it: a steam ID to include in the list.
	Include string `query:"it"`
	// ot: an epoch timestamp. Items older than this timestamp are filtered out.
	ExcludeOlderThan int `query:"ot"`
	// nt: an epoch timestamp. Items newer than this timestamp are filtered out.
	ExcludeNewerThan int `query:"nt"`
}

func (c *RainierController) StreamContents(w http.ResponseWriter, req *http.Request) {
	params := StreamContentsRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamId := chi.URLParam(req, "*"); streamId {
	case "user/-/state/com.google/reading-list":
		c.Container.Render.Text(w, http.StatusOK, "reading list")
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}
