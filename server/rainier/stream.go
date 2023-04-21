package rainier

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// r: sort criteria. Items are sorted by date (descending by default), r=o inverts the order.
// n: the number of items per page. Default: 20.
// c: the continuation string (see below).
// xt: a stream ID to exclude from the list.
// it: a steam ID to include in the list.
// ot: an epoch timestamp. Items older than this timestamp are filtered out.
// nt: an epoch timestamp. Items newer than this timestamp are filtered out.

func (c *RainierController) StreamContents(w http.ResponseWriter, req *http.Request) {
	switch streamID := chi.URLParam(req, "*"); streamID {
	case "user/-/state/com.google/reading-list":
		c.Container.Render.Text(w, http.StatusOK, "reading list")
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}
