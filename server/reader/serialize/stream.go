package serialize

import (
	"github.com/samber/lo"
)

type StreamContentsResponse struct {
	Title        string     `json:"title"`
	Updated      int64      `json:"updated"`
	Continuation string     `json:"continuation"`
	ID           string     `json:"id"`
	Self         Self       `json:"self"`
	Items        []FeedItem `json:"items"`
}

type ReadingListParams struct {
	ID           string
	Title        string
	Continuation string
	SelfURL      string
}

func ReadingList(params ReadingListParams, items []SerializableItem) StreamContentsResponse {
	serialized := FeedItemsFromSerializable(items)

	mostRecent := lo.MaxBy(serialized, func(item FeedItem, max FeedItem) bool {
		return item.Published > max.Published
	})

	return StreamContentsResponse{
		Title:        params.Title,
		Updated:      mostRecent.Published,
		Continuation: "page2", // TODO: paginate
		ID:           params.ID,
		Items:        serialized,
		Self:         Self{Href: params.SelfURL},
	}
}
