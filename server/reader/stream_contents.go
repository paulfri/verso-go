package reader

import (
	"fmt"
	"net/http"
	"net/url"

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

const DEFAULT_ITEMS_PER_PAGE = 1000

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

	rawStreamID := chi.URLParam(req, "*")
	streamID, err := url.QueryUnescape(rawStreamID)
	if err != nil {
		panic(err)
	}

	switch streamIDType := common.StreamIDType(streamID); streamIDType {
	case common.StreamIDReadingList:
		items, err := c.Container.Queries.GetItemsWithURLByUserID(
			ctx,
			query.GetItemsWithURLByUserIDParams{
				UserID: userID,
				Limit:  DEFAULT_ITEMS_PER_PAGE,
			},
		)

		if err != nil {
			panic(err)
		}

		serializable := serialize.QueryRowsToSerializableItems(items)

		response := serialize.ReadingList(
			serialize.ReadingListParams{
				Title: fmt.Sprintf("%s's feed", user.Name),
				ID:    fmt.Sprintf("user/%d/state/com.google/reading-list", user.ID),
				SelfURL: fmt.Sprintf(
					"%s/reader/api/0/stream/contents/user/-/state/com.google/reading-list?output=json", // TODO
					c.Container.Config.BaseURL,
				),
				Continuation: "page2", // TODO: paginate
			},
			serializable,
		)

		c.Container.Render.JSON(w, http.StatusOK, response)
	case common.StreamIDBroadcastFriends:
		// Not implemented.
		c.Container.Render.Text(w, http.StatusNotFound, "")
	case common.StreamIDFormatFeed:
		items, err := c.Container.Queries.GetRecentItemsByRSSFeedURL(
			ctx,
			query.GetRecentItemsByRSSFeedURLParams{
				URL:   common.FeedURLFromReaderStreamID(streamID),
				Limit: DEFAULT_ITEMS_PER_PAGE,
			},
		)

		if err != nil {
			panic(err)
		}

		serializable := serialize.QueryRowsToSerializableItems(items)

		response := serialize.ReadingList(
			serialize.ReadingListParams{
				Title: streamID, // TODO: name of feed?
				ID:    streamID,
				SelfURL: fmt.Sprintf(
					"%s/reader/api/0/stream/contents/%s?output=json",
					c.Container.Config.BaseURL,
					streamID,
				),
				Continuation: "page2", // TODO: paginate
			},
			serializable,
		)

		c.Container.Render.JSON(w, http.StatusOK, response)
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}
