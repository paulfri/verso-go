package reader

import (
	"net/http"

	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"github.com/versolabs/verso/server/reader/serialize"
)

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

	// TODO: we only support one stream ID for now.
	switch streamIDType := common.StreamIDType(params.StreamID); streamIDType {
	case common.StreamIDReadingList:
		items, err := c.Container.Queries.GetQueueItemsByUserID(
			ctx,
			query.GetQueueItemsByUserIDParams{
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
	case common.StreamIDStarred:
		// TODO: Not implemented.
		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: []serialize.FeedItemRef{},
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
	case common.StreamIDRead:
		// TODO: Not implemented.
		c.Container.Render.JSON(w, http.StatusOK, StreamItemsIDsResponse{
			ItemRefs: []serialize.FeedItemRef{},
		})
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}
