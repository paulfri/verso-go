package reader

import (
	"net/http"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"github.com/versolabs/verso/server/reader/serialize"
)

type StreamItemsIDsRequestParams struct {
	StreamID   string `query:"s"`
	ExcludeTag string `query:"xt"`
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
	err := c.Container.BodyOrQueryParams(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamIDType := common.StreamIDType(params.StreamID); streamIDType {
	case common.StreamIDReadingList:
		var itemRefs []serialize.FeedItemRef

		// TODO: only implements exclude tag for Read
		if params.ExcludeTag == common.StreamIDRead {
			itemRefs = c.getUnreadItemsByUserID(req, userID)
		} else {
			itemRefs = c.getAllItemsByUserID(req, userID)
		}

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDRead:
		itemRefs := c.getReadItemsByUserID(req, userID)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDStarred:
		itemRefs := c.getStarredItemsByUserID(req, userID)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDBroadcastFriends:
		// TODO: Not implemented.
		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: []serialize.FeedItemRef{},
		})
	case common.StreamIDFormatFeed:
		// TODO: Not implemented.
		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: []serialize.FeedItemRef{},
		})
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}

func (c *ReaderController) getAllItemsByUserID(req *http.Request, userID int64) []serialize.FeedItemRef {
	ctx := req.Context()
	queries := c.Container.GetQueries(req)

	items, err := queries.GetItemsByUserID(
		ctx,
		query.GetItemsByUserIDParams{
			UserID: userID,
			Limit:  DEFAULT_ITEMS_PER_PAGE,
		},
	)

	if err != nil {
		panic(err) // TODO
	}

	serializable := serialize.QueryRowsToSerializableItems(items)
	itemRefs := serialize.FeedItemRefsFromRows(serializable)

	return itemRefs
}

func (c *ReaderController) getReadItemsByUserID(req *http.Request, userID int64) []serialize.FeedItemRef {
	ctx := req.Context()
	queries := c.Container.GetQueries(req)

	items, err := queries.GetReadItemsByUserID(
		ctx,
		query.GetReadItemsByUserIDParams{
			UserID: userID,
			Limit:  DEFAULT_ITEMS_PER_PAGE,
		},
	)

	if err != nil {
		panic(err) // TODO
	}

	serializable := serialize.QueryRowsToSerializableItems(items)
	itemRefs := serialize.FeedItemRefsFromRows(serializable)

	return itemRefs
}

func (c *ReaderController) getUnreadItemsByUserID(req *http.Request, userID int64) []serialize.FeedItemRef {
	ctx := req.Context()
	queries := c.Container.GetQueries(req)

	items, err := queries.GetUnreadItemsByUserID(
		ctx,
		query.GetUnreadItemsByUserIDParams{
			UserID: userID,
			Limit:  DEFAULT_ITEMS_PER_PAGE,
		},
	)

	if err != nil {
		panic(err) // TODO
	}

	serializable := serialize.QueryRowsToSerializableItems(items)
	itemRefs := serialize.FeedItemRefsFromRows(serializable)

	return itemRefs
}

func (c *ReaderController) getStarredItemsByUserID(req *http.Request, userID int64) []serialize.FeedItemRef {
	ctx := req.Context()
	queries := c.Container.GetQueries(req)

	items, err := queries.GetStarredItemsByUserID(
		ctx,
		query.GetStarredItemsByUserIDParams{
			UserID: userID,
			Limit:  DEFAULT_ITEMS_PER_PAGE,
		},
	)

	if err != nil {
		panic(err) // TODO
	}

	serializable := serialize.QueryRowsToSerializableItems(items)
	itemRefs := serialize.FeedItemRefsFromRows(serializable)

	return itemRefs
}
