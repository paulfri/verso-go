package rainier

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/citra/db/query"
	"gopkg.in/guregu/null.v4"
)

type SortOrderValue string

const (
	Ascending  SortOrderValue = "o"
	Descending SortOrderValue = "a"
)

const DEFAULT_ITEMS_PER_PAGE = 20

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

type StreamContentsResponse struct {
	Direction    string `json:"direction"`
	Author       string `json:"author"`
	Title        string `json:"title"`
	Updated      int64  `json:"updated"`
	Continuation string `json:"continuation"`
	ID           string `json:"id"`
	// TODO
	// "self": [{
	//   "href": "https://api.verso.so/reader/api/0/stream/contents/user/-/state/com.google/reading-list?output=json"
	// }],
	Items []StreamContentsItem `json:"items"`
}

type StreamContentsItemContent struct {
	Direction string `json:"direction"`
	Content   string `json:"content"`
}
type StreamContentsItem struct {
	// TODO:
	//   categories
	//   alternate
	Origin        interface{}               `json:"origin"`
	Updated       int64                     `json:"updated"`
	ID            string                    `json:"id"`
	Author        null.String               `json:"author"`
	Content       StreamContentsItemContent `json:"content"`
	TimestampUsec int64                     `json:"timestampUsec"`
	CrawlTimeMsec int64                     `json:"crawlTimeMsec"`
	Published     int64                     `json:"published"`
	Title         string                    `json:"title"`
}

func (c *RainierController) StreamContents(w http.ResponseWriter, req *http.Request) {
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

	switch streamID := chi.URLParam(req, "*"); streamID {
	case "user/-/state/com.google/reading-list":
		items, err := c.Container.Queries.GetQueueItemsByUserID(
			ctx,
			query.GetQueueItemsByUserIDParams{
				UserID: userID,
				Limit:  DEFAULT_ITEMS_PER_PAGE,
			},
		)

		if err != nil {
			c.Container.Render.Text(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := readingList(user, items)

		c.Container.Render.JSON(w, http.StatusOK, response)
	default:
		c.Container.Render.Text(w, http.StatusBadRequest, "not a stream")
	}
}

func readingList(user query.IdentityUser, items []query.GetQueueItemsByUserIDRow) StreamContentsResponse {
	serialized := lop.Map(items, func(item query.GetQueueItemsByUserIDRow, _ int) StreamContentsItem {
		published := item.CreatedAt.Unix()
		if item.PublishedAt.Valid {
			published = item.PublishedAt.Time.Unix()
		}

		updated := item.CreatedAt.Unix()
		if item.RemoteUpdatedAt.Valid {
			updated = item.RemoteUpdatedAt.Time.Unix()
		}

		return StreamContentsItem{
			Origin: Origin{},
			// TODO: Determine compatibility with Reader
			// https://raw.githubusercontent.com/mihaip/google-reader-api/master/wiki/ItemId.wiki
			ID:     fmt.Sprintf("tag:google.com,2005:reader/item/%016d", item.ID),
			Author: null.String{item.Author},
			Content: StreamContentsItemContent{
				Direction: "ltr",
				Content:   item.Content,
			},
			TimestampUsec: published * 10_000,
			CrawlTimeMsec: item.CreatedAt.Unix() * 1000,
			Published:     published,
			Updated:       updated,
			Title:         item.Title,
		}

	})

	updated := lo.MaxBy(items, func(item query.GetQueueItemsByUserIDRow, max query.GetQueueItemsByUserIDRow) bool {
		return item.PublishedAt.Time.Unix() > max.PublishedAt.Time.Unix()
	})

	return StreamContentsResponse{
		Direction:    "rtl",
		Author:       user.Name,
		Title:        fmt.Sprintf("%s's feed", user.Name),
		Updated:      updated.PublishedAt.Time.Unix(),
		Continuation: "page2", // TODO: paginate
		ID:           fmt.Sprintf("user/%d/state/com.google/reading-list", user.ID),
		Items:        serialized,
	}
}

type Origin struct{}
