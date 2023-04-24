package serialize

import (
	"fmt"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
)

type Subscription struct {
	Title         string `json:"title"`
	FirstItemMsec string `json:"firstitemmsec"`
	HTMLURL       string `json:"htmlUrl"`
	SortID        string `json:"sortid"`
	ID            string `json:"id"`
	// Categories    []struct {
	// 	ID    string `json:"id"`
	// 	Label string `json:"label"`
	// } `json:"categories"`
}

func SubscriptionsFromRows(rows []query.GetSubscriptionsByUserIDRow) []Subscription {
	return lop.Map(rows, func(row query.GetSubscriptionsByUserIDRow, _ int) Subscription {
		return Subscription{
			Title:         row.Title,
			FirstItemMsec: "0", // TODO
			HTMLURL:       row.URL,
			SortID:        fmt.Sprintf("B%07d", row.ID),
			ID:            fmt.Sprintf("%d", row.ID),
		}
	})
}
