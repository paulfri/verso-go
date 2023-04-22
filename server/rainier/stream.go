package rainier

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/citra/db/query"
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
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := StreamContentsRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamID := chi.URLParam(req, "*"); streamID {
	case "user/-/state/com.google/reading-list":
		items, err := c.Container.Queries.GetQueueItemsByUserID(
			ctx,
			query.GetQueueItemsByUserIDParams{
				UserID: userID,
				Limit:  10,
			},
		)

		if err != nil {
			c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		c.Container.Render.JSON(w, http.StatusOK, items)
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}
