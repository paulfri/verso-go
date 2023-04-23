package command

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/versolabs/verso/core/helper"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/worker/tasks"
)

// Subscribes the given user to the given feed URL, creating the feed in the
// database if necessary.
func (c Command) SubscribeToFeedByURL(ctx context.Context, url string, userID int64) error {
	// Track whether we need to scrape the feed to collect metadata.
	needsScrape := false

	// Normalize the given feed URL.
	feedURL, err := helper.NormalizeFeedURL(url)
	if err != nil {
		return err
	}

	// Check if this feed already exists in the database.
	feed, err := c.Queries.GetRSSFeedByURL(ctx, feedURL)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// If this feed doesn't exist, try to find a valid feed at the given
			// normalized URL. We do the lookup first to avoid this network call
			// if we already know that a valid feed exist(ed) at that address.
			feedURL, err = helper.GatherFeed(feedURL)
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
		parsedFeed, err := parser.ParseURL(feedURL)

		if err != nil {
			return err
		}

		title = parsedFeed.Title
	}

	err = c.QueryTransaction(func(withTx *query.Queries) error {
		// Fetch the feed record.
		feed, err := withTx.GetOrCreateRSSFeed(
			ctx,
			query.GetOrCreateRSSFeedParams{
				URL:   feedURL,
				Title: title,
			},
		)

		if err != nil {
			return fmt.Errorf("get or create feed: %v", err)
		}

		// Create a subscription for the user to the feed.
		_, err = withTx.CreateRSSSubscription(
			ctx,
			query.CreateRSSSubscriptionParams{
				UserID: userID,
				FeedID: feed.ID,
			},
		)

		if err != nil {
			return fmt.Errorf("create subscription: %v", err)
		}

		// Fetch the most recent previous 10 items from the RSS feed to the user's queue.
		items, err := withTx.GetRecentItemsByRSSFeedID(
			ctx,
			query.GetRecentItemsByRSSFeedIDParams{
				FeedID: feed.ID,
				Limit:  10,
			},
		)

		if err != nil {
			return fmt.Errorf("get items: %v", err)
		}

		// TODO: Bulk insert items.
		for _, item := range items {
			_, err := withTx.CreateQueueItem(ctx, query.CreateQueueItemParams{
				UserID:    userID,
				RSSItemID: sql.NullInt64{Int64: item.ID, Valid: true},
			})

			// Ignore if the user already has this feed's items in their queue.
			if err != nil && err != sql.ErrNoRows {
				return err
			}
		}

		// Enqueue the feed for parsing.
		task, err := tasks.NewFeedParseTask(feed.ID)
		if err != nil {
			return err
		}

		defer c.Asynq.Enqueue(task)

		return nil
	})

	return err
}
