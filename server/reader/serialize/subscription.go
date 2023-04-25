package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
)

type Subscription struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Categories    []Category `json:"categories"`
	SortID        string     `json:"sortid"`
	FirstItemMsec string     `json:"firstitemmsec"`
	URL           string     `json:"url"`
	HTMLURL       string     `json:"htmlUrl"`
	IconURL       string     `json:"iconUrl"`
}

func SubscriptionsFromRows(rows []query.GetSubscriptionsByUserIDRow) []Subscription {
	return lop.Map(rows, func(row query.GetSubscriptionsByUserIDRow, _ int) Subscription {
		return Subscription{
			ID:            common.ReaderStreamIDFromFeedURL(row.URL),
			Title:         row.Title,
			Categories:    []Category{},
			FirstItemMsec: "123123123", // TODO
			HTMLURL:       row.URL,
			SortID:        fmt.Sprintf("B%07d", row.ID),
			URL:           row.URL,
			IconURL:       "",
		}
	})
}
