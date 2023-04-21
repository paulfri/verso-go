package rainier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/feed"
)

func (c *RainierController) SubscriptionQuickAdd(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	userID := ctx.Value(ContextUserIDKey{}).(int64)
	quickadd := req.URL.Query().Get("quickadd")

	if quickadd == "" {
		w.WriteHeader(400) // TODO: error message
		return
	}

	if !strings.HasPrefix(quickadd, "http://") && !strings.HasPrefix(quickadd, "https://") {
		quickadd = "https://" + quickadd
	}

	url, _ := url.Parse(quickadd)
	normalized := purell.MustNormalizeURLString(url.String(), purell.FlagsSafe)
	parsedFeed, err := feed.Parse(normalized)

	if err == nil {
		tx, err := c.Container.DB.Begin()
		if err != nil {
			panic(err)
		}
		defer tx.Rollback()

		withTx := c.Container.Queries.WithTx(tx)

		feed, err := withTx.FindOrCreateRssFeed(
			ctx,
			query.FindOrCreateRssFeedParams{
				Url:   normalized,
				Title: parsedFeed.Title,
			},
		)

		if err != nil {
			panic(err)
		}

		subscription, err := withTx.CreateSubscription(
			ctx,
			query.CreateSubscriptionParams{
				UserID:    userID,
				RssFeedID: feed.ID,
			},
		)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(feed.Title)
		fmt.Println(subscription.ID)

		err = tx.Commit()
		if err != nil {
			panic(err)
		}

		c.Container.Render.Text(w, http.StatusOK, "ok")
	} else {
		c.Container.Render.Text(w, http.StatusOK, err.Error())
	}
}
