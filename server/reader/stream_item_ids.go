package reader

import (
	"context"
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
	err := c.Container.Params(&params, req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	switch streamIDType := common.StreamIDType(params.StreamID); streamIDType {
	case common.StreamIDReadingList:
		var itemRefs []serialize.FeedItemRef

		// TODO: only implements exclude tag for Read
		if params.ExcludeTag == common.StreamIDRead {
			itemRefs = c.getUnreadItemsByUserID(ctx, userID)
		} else {
			itemRefs = c.getAllItemsByUserID(ctx, userID)
		}

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDRead:
		itemRefs := c.getReadItemsByUserID(ctx, userID)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDStarred:
		itemRefs := c.getStarredItemsByUserID(ctx, userID)

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

func (c *ReaderController) getAllItemsByUserID(ctx context.Context, userID int64) []serialize.FeedItemRef {
	items, err := c.Container.Queries.GetItemsByUserID(
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

func (c *ReaderController) getReadItemsByUserID(ctx context.Context, userID int64) []serialize.FeedItemRef {
	items, err := c.Container.Queries.GetReadItemsByUserID(
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

func (c *ReaderController) getUnreadItemsByUserID(ctx context.Context, userID int64) []serialize.FeedItemRef {
	items, err := c.Container.Queries.GetUnreadItemsByUserID(
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

func (c *ReaderController) getStarredItemsByUserID(ctx context.Context, userID int64) []serialize.FeedItemRef {
	items, err := c.Container.Queries.GetStarredItemsByUserID(
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

// Don't know how to begin to explain this
