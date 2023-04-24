package reader

import (
	"fmt"
	"net/http"

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

	streamID := chi.URLParam(req, "*")

	switch streamIDType := common.StreamIDType(streamID); streamIDType {
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
	case common.StreamIDFormatFeed:
		// TODO
		c.Container.Render.Text(w, http.StatusBadRequest, "not yet implemented")
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}

type StreamItemsIDsRequestParams struct {
	StreamID string `query:"s"`
	// TODO implement
	NumItems int `query:"n"`
}

type StreamItemsIDsResponse struct {
	ItemRefs []serialize.FeedItemRef `json:"itemRefs"`
}

func (c *ReaderController) StreamItemsIDs(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := StreamItemsIDsRequestParams{}
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamIDType := common.StreamIDType(params.StreamID); streamIDType {
	case common.StreamIDReadingList:
		fallthrough
	case common.StreamIDStarred:
		items, err := c.Container.Queries.GetQueueItemsByUserID(
			ctx,
			query.GetQueueItemsByUserIDParams{
				UserID: userID,
				Limit:  DEFAULT_ITEMS_PER_PAGE,
			},
		)

		if err != nil {
			panic(err)
		}

		itemRefs := serialize.FeedItemRefsFromRows(items)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{ItemRefs: itemRefs})
	case common.StreamIDBroadcastFriends:
		// Not implemented.
		c.Container.Render.Text(w, http.StatusNotFound, "")
	case common.StreamIDFormatFeed:
		// TODO
		c.Container.Render.Text(w, http.StatusBadRequest, "not yet implemented")
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}

type StreamItemsContentsRequestParams struct {
	Items []string `query:"i"`
	// || `output` || Output format. `json`, `atom` RFC 4287, `atom-hifi`
	// || `sharers` // || `likes` // || `comments` // || `trans` // || `mediaRss`
}

func (c *ReaderController) StreamItemsContents(w http.ResponseWriter, req *http.Request) {
	// ctx := req.Context()
	// userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := StreamItemsContentsRequestParams{}
	err := c.Container.BodyParams(&params, req)

	fmt.Printf("%+v\n", params)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	// switch streamIDType := common.StreamIDType(params.StreamID); streamIDType {
	// case common.StreamIDReadingList:
	// 	fallthrough
	// case common.StreamIDStarred:
	// 	items, err := c.Container.Queries.GetQueueItemsByUserID(
	// 		ctx,
	// 		query.GetQueueItemsByUserIDParams{
	// 			UserID: userID,
	// 			Limit:  DEFAULT_ITEMS_PER_PAGE,
	// 		},
	// 	)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	itemRefs := serialize.FeedItemRefsFromRows(items)

	// 	c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{ItemRefs: itemRefs})
	// case common.StreamIDBroadcastFriends:
	// 	// Not implemented.
	// 	c.Container.Render.Text(w, http.StatusNotFound, "")
	// case common.StreamIDFormatFeed:
	// 	// TODO
	// 	c.Container.Render.Text(w, http.StatusBadRequest, "not yet implemented")
	// default:
	// 	c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	// }
}
