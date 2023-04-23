package serialize

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/versolabs/citra/db/query"
)

type StreamContentsResponse struct {
	Direction    string     `json:"direction"`
	Author       string     `json:"author"`
	Title        string     `json:"title"`
	Updated      int64      `json:"updated"`
	Continuation string     `json:"continuation"`
	ID           string     `json:"id"`
	Self         Self       `json:"self"`
	Items        []FeedItem `json:"items"`
}

func ReadingList(user query.IdentityUser, items []query.GetQueueItemsByUserIDRow, baseURL string) StreamContentsResponse {
	serialized := FeedItemsFromRows(items)

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
		Self: Self{
			Href: fmt.Sprintf(
				"%s/reader/api/0/stream/contents/user/-/state/com.google/reading-list?output=json",
				baseURL,
			),
		},
	}
}
