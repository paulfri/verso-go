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

func (c *Controller) TagList(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(ContextUserIDKey{}).(int64)
	queries := c.Container.GetQueries(req)

	tagRows, err := queries.GetTagsByUserID(req.Context(), userID)

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

func (c *Controller) EditTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := EditTagRequestParams{}
	err := c.Container.BodyOrQueryParams(&params, req)

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

func (c *Controller) addTag(ctx context.Context, readerID string, userID int64, tag string) error {
	switch tag {
	case common.StreamIDRead:
		return c.Container.Command.MarkRead(ctx, readerID, userID)
	case common.StreamIDStarred:
		return c.Container.Command.MarkStarred(ctx, readerID, userID)
	default:
		return fmt.Errorf("unknown tag: %s", tag)
	}
}

func (c *Controller) removeTag(ctx context.Context, readerID string, userID int64, tag string) error {
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

func (c *Controller) DisableTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := DisableTagRequestParams{}
	err := c.Container.BodyOrQueryParams(&params, req)
	queries := c.Container.GetQueries(req)

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

	err = queries.DeleteTagByNameAndUserID(ctx, query.DeleteTagByNameAndUserIDParams{
		Name:   tagWithoutPrefix,
		UserID: userID,
	})

	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}

type RenameTagRequestParams struct {
	Tag     string `query:"s" validate:"required"`
	NewName string `query:"dest" validate:"required"`
}

func (c *Controller) RenameTag(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	params := RenameTagRequestParams{}
	queries := c.Container.GetQueries(req)
	err := c.Container.BodyOrQueryParams(&params, req)

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
	newTagWithoutPrefix, _ := strings.CutPrefix(params.NewName, "user/-/label/")

	err = queries.RenameTagByNameAndUserID(ctx, query.RenameTagByNameAndUserIDParams{
		UserID:  userID,
		Name:    tagWithoutPrefix,
		NewName: newTagWithoutPrefix,
	})

	if err != nil && err != sql.ErrNoRows {
		panic(err) // TODO
	}

	c.Container.Render.Text(w, http.StatusOK, "OK")
}
