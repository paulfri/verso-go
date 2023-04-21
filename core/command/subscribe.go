package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/versolabs/citra/core/helper"
	"github.com/versolabs/citra/db/query"
)

// Subscribes the given user to the given feed URL, creating the feed in
// the database if necessary.
func (c Command) SubscribeToFeedByUrl(ctx context.Context, url string, userId int64) error {
	// Track whether we need to scrape the feed to collect metadata.
	needsScrape := false

	// Normalize the given feed URL.
	feedUrl, err := helper.NormalizeFeedUrl(url)
	if err != nil {
		return err
	}

	// Check if this feed already exists in the database.
	feed, err := c.Queries.FindRssFeedByUrl(ctx, feedUrl)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// If this feed doesn't exist, try to find a valid feed at the given
			// normalized URL. We do the lookup first to avoid this network call
			// if we already know that a valid feed exist(ed) at that address.
			feedUrl, err = helper.GatherFeed(feedUrl)
			// We'll need the title by scraping the feed.
			needsScrape = true

			if err != nil {
				return err
			}
		default:
			// If this was a general database error, return that instead of
			// continuing.
			return err
		}
	}

	// If we need to scrape the feed to get its metadata (i.e. it's a new feed),
	// then do that here.
	title := feed.Title
	if needsScrape {
		// TOOD: this can probably be a helper
		parser := gofeed.NewParser()
		parsedFeed, err := parser.ParseURL(feedUrl)

		if err != nil {
			return err
		}

		title = parsedFeed.Title
	}

	err = c.QueryTransaction(func(withTx *query.Queries) error {
		feed, err := withTx.FindOrCreateRssFeed(
			ctx,
			query.FindOrCreateRssFeedParams{
				Url:   feedUrl,
				Title: title,
			},
		)

		if err != nil {
			return fmt.Errorf("find or create feed: %v", err)
		}

		_, err = withTx.CreateSubscription(
			ctx,
			query.CreateSubscriptionParams{
				UserID:    userId,
				RssFeedID: feed.ID,
			},
		)

		if err != nil {
			return fmt.Errorf("create subscription: %v", err)
		}

		return nil
	})

	return err
}
