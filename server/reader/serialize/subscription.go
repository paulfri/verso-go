package serialize

import (
	"fmt"

	"github.com/samber/lo"
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
	// Group the rows by subscription ID.
	rowsByID := lo.GroupBy(rows, func(i query.GetSubscriptionsByUserIDRow) int64 {
		return i.FeedID
	})

	return lo.MapToSlice(rowsByID, func(id int64, tagRows []query.GetSubscriptionsByUserIDRow) Subscription {
		categories := lo.FilterMap(tagRows, func(row query.GetSubscriptionsByUserIDRow, _ int) (Category, bool) {
			if !row.Name.Valid {
				return Category{}, false
			}

			return Category{
				ID:    fmt.Sprintf("user/-/label/%s", row.Name.String),
				Label: row.Name.String,
			}, true
		})

		sub := tagRows[0]
		url := sub.RSSFeedURL

		return Subscription{
			ID:            common.ReaderStreamIDFromFeedURL(url),
			Title:         sub.Title,
			Categories:    categories,
			FirstItemMsec: "123123123",                  // TODO
			SortID:        fmt.Sprintf("B%07d", sub.ID), // TODO
			HTMLURL:       url,                          // TODO
			URL:           url,                          // TODO
			IconURL:       "",                           // TODO
		}
	})
}
