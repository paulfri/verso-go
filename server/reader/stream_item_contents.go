package reader

import (
	"fmt"
	"net/http"
	"time"

	"github.com/samber/lo"
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

func (c *ReaderController) StreamItemsContents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := StreamItemsContentsRequestParams{}
	err := c.Container.BodyParams(&params, req)

	fmt.Printf("%+v\n", params)

	if err != nil {
		c.Container.Render.Text(w, http.StatusBadRequest, err.Error())
		return
	}

	readerIDs := lo.Map(params.Items, func(itemID string, _ int) int64 {
		return common.ReaderIDFromInput(itemID)
	})

	items, err := c.Container.Queries.GetQueueItemsByReaderIDs(ctx, readerIDs)
	if err != nil {
		panic(err) // TODO
	}

	serialized := serialize.FeedItemsFromReaderIDsRows(items)

	var id string
	if len(items) > 0 {
		id = common.LongItemID(items[0].ID)
	} else {
		id = ""
	}

	c.Container.Render.JSON(w, http.StatusOK, serialize.StreamContentsResponse{
		ID:           id,
		Title:        "asdf",
		Continuation: "TOOD",
		Self: serialize.Self{
			Href: "http://localhost:8080/reader/api/0/stream/items/contents", // TODO
		},
		Updated: time.Now().Unix(),
		Items:   serialized,
	})
}
