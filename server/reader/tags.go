package reader

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
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
	userID := req.Context().Value(ContextUserIDKey{}).(int64)

	tagRows, err := c.Container.Queries.GetTagsByUserID(req.Context(), userID)

	if err != nil {
		panic(err)
	}

	var tags []Tag
	tags = append(tags, Tag{
		ID:     "user/-/state/com.google/starred",
		SortID: "A0000001",
	})

	userTags := lop.Map(tagRows, func(row query.TaxonomyTag, index int) Tag {
		return Tag{
			ID:     fmt.Sprintf("user/-/label/%s", row.Name),
			SortID: fmt.Sprintf("A%07d", index+len(tags)+1),
		}
	})

	tags = append(tags, userTags...)

	c.Container.Render.JSON(w, http.StatusOK, TagList{
		Tags: tags,
	})
}

type EditTagRequestParams struct {
	ItemIDs   []string `query:"i" validate:"required"`
	AddTag    string   `query:"a" validate:"required_without=RemoveTag"`
	RemoveTag string   `query:"r" validate:"required_without=AddTag"`
}

func (c *ReaderController) EditTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := EditTagRequestParams{}
	err := c.Container.BodyParams(&params, req)

	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Container.Validator.Struct(params)
	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: batch this
	for _, item := range params.ItemIDs {
		readerID := common.ReaderIDFromInput(item)

		if params.AddTag != "" {
			err = c.addTag(ctx, readerID, userID, params.AddTag)
		} else {
			err = c.removeTag(ctx, readerID, userID, params.RemoveTag)
		}
	}

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}

func (c *ReaderController) addTag(ctx context.Context, readerID string, userID int64, tag string) error {
	switch tag {
	case common.StreamIDRead:
		return c.Container.Command.MarkRead(ctx, readerID, userID)
	case common.StreamIDStarred:
		return c.Container.Command.MarkStarred(ctx, readerID, userID)
	default:
		return fmt.Errorf("unknown tag: %s", tag)
	}
}

func (c *ReaderController) removeTag(ctx context.Context, readerID string, userID int64, tag string) error {
	switch tag {
	case common.StreamIDRead:
		return c.Container.Command.MarkUnread(ctx, readerID, userID)
	case common.StreamIDStarred:
		return c.Container.Command.MarkUnstarred(ctx, readerID, userID)
	default:
		return fmt.Errorf("unknown tag: %s", tag)
	}
}

type DisableTagRequestParams struct {
	Tag string `query:"s" validate:"required"`
}

func (c *ReaderController) DisableTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := DisableTagRequestParams{}
	err := c.Container.BodyParams(&params, req)

	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Container.Validator.Struct(params)
	if err != nil {
		c.Container.Render.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	tagWithoutPrefix, _ := strings.CutPrefix(params.Tag, "user/-/label/")

	err = c.Container.Queries.DeleteTagByNameAndUserID(ctx, query.DeleteTagByNameAndUserIDParams{
		Name:   tagWithoutPrefix,
		UserID: userID,
	})

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}
