package reader

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
)

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
