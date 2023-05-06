package reader

import (
	"net/http"
	"time"

	"github.com/samber/lo"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
	"github.com/versolabs/verso/server/reader/serialize"
)

type StreamItemsContentsRequestParams struct {
	Items []string `query:"i"`
	// Output string `json:"output"` // json, atom, atom-hifi
	// Sharers `json:"sharers"`
	// Likes `json:"likes"`
	// Comments `json:"comments"`
	// Trans `json:"trans"`
	// MediaRSS `json:"mediaRss"`
}

func (c *Controller) StreamItemsContents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := StreamItemsContentsRequestParams{}
	err := c.Container.BodyOrQueryParams(&params, req)
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	queries := c.Container.GetQueries(req)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())

		return
	}

	readerIDs := lo.Map(params.Items, func(itemID string, _ int) string {
		return common.ReaderIDFromInput(itemID)
	})

	items, err := queries.GetItemsWithContentDataByReaderIDs(
		ctx,
		query.GetItemsWithContentDataByReaderIDsParams{
			Column1: readerIDs,
			UserID:  userID,
		},
	)
	if err != nil {
		panic(err) // TODO
	}

	serializable := serialize.QueryRowsToSerializableItems(items)
	serialized := serialize.FeedItemsFromSerializable(serializable)

	var id string
	if len(items) > 0 {
		id = common.LongItemID(items[0].ReaderID)
	} else {
		id = ""
	}

	c.Container.Render.JSON(w, http.StatusOK, serialize.StreamContentsResponse{
		ID:           id,
		Title:        "asdf",
		Continuation: "TODO",
		Self: serialize.Self{
			Href: "http://localhost:8080/reader/api/0/stream/items/contents", // TODO
		},
		Updated: time.Now().Unix(),
		Items:   serialized,
	})
}
