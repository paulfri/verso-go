package reader

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"gopkg.in/guregu/null.v4"
)

type SubscriptionQuickAddRequestParams struct {
	Quickadd string `query:"quickadd" validate:"required"`
}

type SubscriptionQuickAddResponse struct {
	NumResults int8        `json:"numResults"`
	Query      string      `json:"query"`
	StreamID   null.String `json:"streamId"`
}

func (c *ReaderController) SubscriptionQuickAdd(w http.ResponseWriter, req *http.Request) {
	params := SubscriptionQuickAddRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)

	subscription, err := c.Container.Command.SubscribeToFeedByURL(ctx, params.Quickadd, userID)
	feed, err2 := c.Container.Queries.GetRSSFeed(ctx, subscription.FeedID)

	if err != nil && err2 != nil {
		c.Container.Render.JSON(
			w,
			http.StatusBadRequest,
			SubscriptionQuickAddResponse{
				NumResults: 0,
				Query:      params.Quickadd,
				StreamID:   null.String{},
			},
		)
	} else {
		streamID := common.ReaderStreamIDFromFeedURL(feed.URL)
		asSQL := sql.NullString{String: streamID, Valid: true}
		asNullable := null.String{asSQL}

		c.Container.Render.JSON(
			w,
			http.StatusOK,
			SubscriptionQuickAddResponse{
				NumResults: 1,
				Query:      params.Quickadd,
				StreamID:   asNullable,
			},
		)
	}
}

type SubscriptionExistsRequestParams struct {
	StreamID string `query:"s" validate:"required"`
}

func (c *ReaderController) SubscriptionExists(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := SubscriptionExistsRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	url := common.FeedURLFromReaderStreamID(params.StreamID)
	feed, err := c.Container.Queries.GetRSSFeedByURL(ctx, url)

	if err == sql.ErrNoRows {
		c.Container.Render.Text(w, http.StatusOK, strconv.FormatBool(false))
		return
	}

	_, err = c.Container.Queries.GetSubscriptionByRSSFeedIDAndUserID(
		ctx,
		query.GetSubscriptionByRSSFeedIDAndUserIDParams{
			FeedID: feed.ID,
			UserID: userID,
		},
	)

	if err == sql.ErrNoRows {
		c.Container.Render.Text(w, http.StatusOK, strconv.FormatBool(false))
		return
	} else if err != nil {
		panic(err)
	}

	c.Container.Render.Text(w, http.StatusOK, strconv.FormatBool(true))
}
