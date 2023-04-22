package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/worker/tasks"
)

func (worker *Worker) HandleFeedParseTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.FeedParsePayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Parsing feed: feed_id=%d\n", p.FeedID)

	thisFeed, err := worker.Container.Queries.GetRSSFeed(ctx, int64(p.FeedID))
	if err != nil {
		return err
	}

	fmt.Println(thisFeed)

	// Fetch the feed and parse it.
	parser := gofeed.NewParser()
	remoteFeed, _ := parser.ParseURL(thisFeed.URL)

	items := remoteFeed.Items
	if len(items) == 0 {
		// Empty feed. Suspicious.
		// TODO: probably log this error somewhere
		return nil
	}

	subscribers, err := worker.Container.Queries.GetSubscribersByRSSFeedID(ctx, thisFeed.ID)
	if err != nil {
		return err
	}

	for _, feedItem := range remoteFeed.Items {
		fmt.Println(feedItem.Title)

		rssItem, err := worker.Container.Queries.CreateRSSItem(ctx, query.CreateRSSItemParams{
			FeedID:      int64(p.FeedID),
			RSSGuid:     feedItem.GUID,
			Title:       feedItem.Title,
			Content:     feedItem.Content,
			Link:        feedItem.Link,
			PublishedAt: sql.NullTime{Time: *feedItem.PublishedParsed, Valid: true},
			// TODO: figure out what's up with this
			// RemoteUpdatedAt: sql.NullTime{Time: *item.UpdatedParsed, Valid: true},
		})

		for _, subscription := range subscribers {
			// TODO: handle error
			worker.Container.Queries.CreateQueueItem(ctx, query.CreateQueueItemParams{
				UserID:    subscription.UserID,
				RSSItemID: sql.NullInt64{Int64: rssItem.ID, Valid: true},
			})
		}

		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
