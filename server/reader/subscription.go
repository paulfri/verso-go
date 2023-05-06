package reader

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"github.com/versolabs/verso/server/reader/serialize"
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

func (c *Controller) SubscriptionQuickAdd(w http.ResponseWriter, req *http.Request) {
	params := SubscriptionQuickAddRequestParams{}
	err := c.Container.BodyOrQueryParams(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())

		return
	}

	ctx := req.Context()
	queries := c.Container.GetQueries(req)
	userID := ctx.Value(ContextUserIDKey{}).(int64)

	subscription, err := c.Container.Command.SubscribeToFeedByURL(ctx, params.Quickadd, userID)

	if err != nil {
		c.Container.Render.JSON(
			w,
			http.StatusBadRequest,
			SubscriptionQuickAddResponse{
				NumResults: 0,
				Query:      params.Quickadd,
				StreamID:   null.String{},
			},
		)

		return
	}

	feed, err2 := queries.GetRSSFeed(ctx, subscription.FeedID)

	if err2 != nil {
		c.Container.Render.JSON(
			w,
			http.StatusBadRequest,
			SubscriptionQuickAddResponse{
				NumResults: 0,
				Query:      params.Quickadd,
				StreamID:   null.String{},
			},
		)

		return
	}

	streamID := common.ReaderStreamIDFromFeedURL(feed.URL)
	asSQL := sql.NullString{String: streamID, Valid: true}

	c.Container.Render.JSON(
		w,
		http.StatusOK,
		SubscriptionQuickAddResponse{
			NumResults: 1,
			Query:      params.Quickadd,
			StreamID:   null.String{NullString: asSQL},
		},
	)
}

type SubscriptionExistsRequestParams struct {
	StreamID string `query:"s" validate:"required"`
}

func (c *Controller) SubscriptionExists(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := SubscriptionExistsRequestParams{}
	err := c.Container.Params(&params, req)
	queries := c.Container.GetQueries(req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())

		return
	}

	url := common.FeedURLFromReaderStreamID(params.StreamID)
	feed, err := queries.GetRSSFeedByURL(ctx, url)

	if err == sql.ErrNoRows {
		c.Container.Render.Text(w, http.StatusOK, strconv.FormatBool(false))

		return
	}

	_, err = queries.GetSubscriptionByRSSFeedIDAndUserID(
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

type SubscriptionListResponse struct {
	Subscriptions []serialize.Subscription `json:"subscriptions"`
}

func (c *Controller) SubscriptionList(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	queries := c.Container.GetQueries(req)
	subscriptions, err := queries.GetSubscriptionsByUserID(ctx, userID)

	if err != nil {
		panic(err) // TODO: fix
	}

	serialized := serialize.SubscriptionsFromRows(subscriptions)

	response := SubscriptionListResponse{
		Subscriptions: serialized,
	}

	c.Container.Render.JSON(w, http.StatusOK, response)
}

type SubscriptionEditParams struct {
	Action    string `query:"ac" validate:"required,eq=subscribe|eq=edit|eq=unsubscribe"`
	StreamID  string `query:"s" validate:"required,startswith=feed/"`
	Title     string `query:"t"`
	AddTag    string `query:"a"`
	RemoveTag string `query:"r"`
}

func (c *Controller) SubscriptionEdit(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := SubscriptionEditParams{}
	err := c.Container.BodyOrQueryParams(&params, req)
	queries := c.Container.GetQueries(req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())

		return
	}

	feedURL := common.FeedURLFromReaderStreamID(params.StreamID)

	switch params.Action {
	case "subscribe":
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error()) // TODO

		return
	case "edit":
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error()) // TODO

		return
	case "unsubscribe":
		queries.DeleteSubscriptionByRSSFeedURLAndUserID(
			ctx,
			query.DeleteSubscriptionByRSSFeedURLAndUserIDParams{
				UserID:     userID,
				RSSFeedURL: feedURL,
			},
		)
	default:
		break
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}
