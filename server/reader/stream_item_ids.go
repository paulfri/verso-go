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
		var items []query.RSSItem

		// TODO: only implements exclude tag for Read
		if params.ExcludeTag == common.StreamIDRead {
			items = c.getUnreadItemsByUserID(ctx, userID)
		} else {
			items = c.getAllItemsByUserID(ctx, userID)
		}

		itemRefs := serialize.FeedItemRefsFromRows(items)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDRead:
		items := c.getReadItemsByUserID(ctx, userID)
		itemRefs := serialize.FeedItemRefsFromRows(items)

		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: itemRefs,
		})
	case common.StreamIDStarred:
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

		itemRefs := serialize.FeedItemRefsFromRows(items)

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

func (c *ReaderController) getAllItemsByUserID(ctx context.Context, userID int64) []query.RSSItem {
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

	return items
}

func (c *ReaderController) getReadItemsByUserID(ctx context.Context, userID int64) []query.RSSItem {
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

	return items
}

func (c *ReaderController) getUnreadItemsByUserID(ctx context.Context, userID int64) []query.RSSItem {
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

	return items
}
