package reader

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/versolabs/verso/server/reader/common"
)

type Tag struct {
	ID     string `json:"id"`
	SortID string `json:"sortid"` // e.g. A0000001
}

type TagList struct {
	Tags []Tag `json:"tags"`
}

func (c *ReaderController) TagList(w http.ResponseWriter, req *http.Request) {
	c.Container.Render.JSON(w, http.StatusOK, TagList{
		Tags: []Tag{},
	})
}

// {
//     "tags": [
//         {
//             "id": "user/1/state/com.google/starred",
//             "sortid": "A0000001"

//         },
//         {
//             "id": "user/1/states/com.google/broadcast",
//             "sortid": "A0000002"

//	        },
//	        {
//	            "id": "user/1/label/Tech",
//	            "sortid": "A0000003"
//	        },
//	    ]
//	}

type EditTagRequestParams struct {
	// This endpoint only accepts one ItemID, but because the middleware modifies
	// the request body to convert `i` to `i[]` to support other endpoints, this
	// struct takes an array.
	ItemID    []string `query:"i" validate:"required"`
	AddTag    string   `query:"a" validate:"required_without=RemoveTag"`
	RemoveTag string   `query:"r" validate:"required_without=AddTag"`
}

func (c *ReaderController) EditTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := EditTagRequestParams{}
	err := c.Container.BodyParams(&params, req)

	fmt.Printf("%+v\n", params)
	fmt.Printf("%+v\n", params)
	fmt.Printf("%+v\n", params)
	fmt.Printf("%+v\n", params)

	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Container.Validator.Struct(params)
	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// See note in EditTagRequestParams.
	item := params.ItemID[0]
	itemID := common.ReaderIDFromInput(item)

	fmt.Println(itemID)
	fmt.Println(itemID)
	fmt.Println(itemID)
	fmt.Println(itemID)
	fmt.Println(itemID)

	if params.AddTag != "" {
		err = c.addTag(ctx, itemID, userID, params.AddTag)
	} else {
		err = c.removeTag(ctx, itemID, userID, params.RemoveTag)
	}

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}

func (c *ReaderController) addTag(ctx context.Context, itemID int64, userID int64, tag string) error {
	switch tag {
	case common.StreamIDRead:
		return c.Container.Command.MarkRead(ctx, itemID, userID)
	case common.StreamIDStarred:
		return c.Container.Command.MarkStarred(ctx, itemID, userID)
	default:
		return fmt.Errorf("unknown tag: %s", tag)
	}
}

func (c *ReaderController) removeTag(ctx context.Context, itemID int64, userID int64, tag string) error {
	switch tag {
	case common.StreamIDRead:
		return c.Container.Command.MarkUnread(ctx, itemID, userID)
	case common.StreamIDStarred:
		return c.Container.Command.MarkUnstarred(ctx, itemID, userID)
	default:
		return fmt.Errorf("unknown tag: %s", tag)
	}
}
