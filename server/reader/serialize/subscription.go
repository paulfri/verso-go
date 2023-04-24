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

// id: "feed/http://www.theanimationblog.com/feed/",
// title: "The Animation Blog.com | Est. 2007",
// categories: [{
// 		id: "user/1005921515/label/Animation",
// 		label: "Animation"
// }],
// sortid: "00DA6134",
// firstitemmsec: 1424501776942006,
// url: "http://www.theanimationblog.com/feed/",
// htmlUrl: "http://www.theanimationblog.com/",
// iconUrl: ""

func SubscriptionsFromRows(rows []query.GetSubscriptionsByUserIDRow) []Subscription {
	return lop.Map(rows, func(row query.GetSubscriptionsByUserIDRow, _ int) Subscription {
		return Subscription{
			ID:            common.ReaderStreamIDFromFeedURL(row.URL),
			Title:         row.Title,
			Categories:    []Category{},
			FirstItemMsec: "123123123", // TODO
			HTMLURL:       row.URL,
			SortID:        fmt.Sprintf("B%07d", row.ID),
			URL:           "",
			IconURL:       "",
		}
	})
}
