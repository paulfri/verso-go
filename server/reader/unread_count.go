package reader

import (
	"fmt"
	"net/http"
	"time"

	lop "github.com/samber/lo/parallel"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/server/reader/common"
)

type UnreadCount struct {
	Count                   int64  `json:"count"`
	ID                      string `json:"id"`
	NewestItemTimestampUSec string `json:"newestItemTimestampUsec"`
}

type UnreadCountResponse struct {
	Max          int           `json:"max"`
	UnreadCounts []UnreadCount `json:"unreadcounts"`
}

func (c *ReaderController) UnreadCount(w http.ResponseWriter, req *http.Request) {
	counts, err := c.Container.Queries.GetUnreadCountsByUserID(req.Context(), 1)

	if err != nil {
		panic(err) // TODO
	}

	unreadCounts := lop.Map(counts, func(row query.GetUnreadCountsByUserIDRow, _ int) UnreadCount {
		newest := row.Newest.(time.Time)
		newestUnixNano := newest.UnixNano()

		return UnreadCount{
			Count:                   row.Count,
			ID:                      common.ReaderStreamIDFromFeedURL(row.URL),
			NewestItemTimestampUSec: fmt.Sprintf("%d", newestUnixNano),
		}
	})

	c.Container.Render.JSON(w, http.StatusOK, UnreadCountResponse{
		Max:          1000, // TODO
		UnreadCounts: unreadCounts,
	})
}

// {
//     "max": 1000,
//     "unreadcounts": [
//         {
//             "count": 1,
//             "id": "feed/http://rss.slashdot.org/Slashdot/slashdot",
//             "newestItemTimestampUsec": "1405452360000000"
//         },
//         {
//             "count": 1,
//             "id": "feed/http://feeds.feedburner.com/alistapart/main",
//             "newestItemTimestampUsec": "1405432727000000"
//         },
//         {
//             "count": 2,
//             "id": "user/1/label/Tech",
//             "newestItemTimestampUsec": "1405432727000000"
//         },
//         {
//             "count": 2,
//             "id": "user/1/state/com.google/reading-list",
//             "newestItemTimestampUsec": "1405432727000000"
//         }
//     ]
// }
