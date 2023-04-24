package reader

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"github.com/versolabs/verso/server/reader/serialize"
)

type SortOrderValue string

const (
	Ascending  SortOrderValue = "o"
	Descending SortOrderValue = "a"
)

const DEFAULT_ITEMS_PER_PAGE = 20

type StreamContentsRequestParams struct {
	// r: sort criteria. Items are sorted by date (descending by default), r=o inverts the order.
	SortOrder SortOrderValue `query:"r"`
	// n: the number of items per page. Default: 20.
	NumItems int `query:"n"`
	// c: the continuation string
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

func (c *ReaderController) StreamContents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := StreamContentsRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.Container.Queries.GetUser(ctx, userID)
	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamID := chi.URLParam(req, "*"); streamID {
	case common.StreamIDReadingList:
		items, err := c.Container.Queries.GetQueueItemsByUserID(
			ctx,
			query.GetQueueItemsByUserIDParams{
				UserID: userID,
				Limit:  DEFAULT_ITEMS_PER_PAGE,
			},
		)

		if err != nil {
			c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := serialize.ReadingList(user, items, c.Container.Config.BaseURL)

		c.Container.Render.JSON(w, http.StatusOK, response)
	case common.StreamIDBroadcastFriends:
		// Not implemented.
		c.Container.Render.Text(w, http.StatusNotFound, "")
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}

type StreamMarkAllAsReadRequestBody struct {
	StreamID  string `json:"s" validate:"required"`
	Timestamp int64  `json:"ts"`
}

func (c *ReaderController) StreamMarkAllAsRead(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	body := StreamMarkAllAsReadRequestBody{}
	err := c.Container.JSONBody(req, &body)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}
	c.Container.Logger.Debug().Msg(fmt.Sprintf("%v", body))

	var publishedBefore sql.NullTime
	if body.Timestamp > 0 {
		// TODO: Is timestamp in milliseconds, microseconds, nanoseconds?
		// Assumed milliseconds here.
		publishedBefore = sql.NullTime{Time: time.Unix(body.Timestamp/1000, 0), Valid: true}
	} else {
		publishedBefore = sql.NullTime{Valid: false}
	}

	switch streamIDType := common.StreamIDType(body.StreamID); streamIDType {
	case common.StreamIDReadingList:
		err := c.Container.Queries.MarkAllQueueItemsAsRead(
			ctx,
			query.MarkAllQueueItemsAsReadParams{
				UserID:          userID,
				RSSFeedURL:      sql.NullString{},
				PublishedBefore: publishedBefore,
			},
		)

		if err != nil {
			c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
			return
		}
	case common.StreamIDFormatFeed:
		feedURL := common.FeedURLFromReaderStreamID(body.StreamID)
		err := c.Container.Queries.MarkAllQueueItemsAsRead(
			ctx,
			query.MarkAllQueueItemsAsReadParams{
				UserID:          userID,
				RSSFeedURL:      feedURL,
				PublishedBefore: publishedBefore,
			},
		)

		if err != nil {
			c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not supported yet")
		return
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}
