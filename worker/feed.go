package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/worker/tasks"
)

func (worker *Worker) HandleFeedParseTask(ctx context.Context, t *asynq.Task) error {
	logger := worker.Container.Logger.With().Str("task", t.Type()).Logger()
	var p tasks.FeedParsePayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	logger.Info().Msgf("Parsing feed: feed_id=%d\n", p.FeedID)

	thisFeed, err := worker.Container.Queries.GetRSSFeed(ctx, int64(p.FeedID))
	if err != nil {
		return err
	}

	// Fetch the feed and parse it.
	parser := gofeed.NewParser()
	remoteFeed, err := parser.ParseURL(thisFeed.URL)

	if err != nil {
		// Parsing error. This will get re-enqueued, so failing is fine.
		// TODO: Store parsing errors for investigation.
		logger.Error().Msgf("Failed to parse feed: %v", err.Error())

		return nil
	}

	items := remoteFeed.Items
	if len(items) == 0 {
		// Empty feed. Suspicious.
		// TODO: probably log this error somewhere

		return nil
	}

	subscribers, err := worker.Container.Queries.GetSubscriptionsByRSSFeedID(ctx, thisFeed.ID)
	if err != nil {
		return err
	}

	for _, item := range items {
		author := sql.NullString{}
		authorEmail := sql.NullString{}

		if len(item.Authors) > 0 {
			a := item.Authors[0]
			if a.Name != "" {
				author = sql.NullString{String: a.Name, Valid: true}
			}
			if a.Email != "" {
				authorEmail = sql.NullString{String: a.Email, Valid: true}
			}
		}

		published := sql.NullTime{}
		if item.PublishedParsed != nil {
			published = sql.NullTime{Time: *item.PublishedParsed, Valid: true}
		}
		updated := sql.NullTime{}
		if item.UpdatedParsed != nil {
			updated = sql.NullTime{Time: *item.UpdatedParsed, Valid: true}
		}

		rssItem, err := worker.Container.Queries.CreateRSSItem(
			ctx,
			query.CreateRSSItemParams{
				FeedID:          int64(p.FeedID),
				RSSGuid:         item.GUID,
				Title:           item.Title,
				Content:         item.Content,
				Author:          author,
				AuthorEmail:     authorEmail,
				Link:            item.Link,
				PublishedAt:     published,
				RemoteUpdatedAt: updated,
			},
		)

		for _, subscription := range subscribers {
			// TODO: handle error
			worker.Container.Queries.CreateQueueItem(ctx, query.CreateQueueItemParams{
				UserID:         subscription.UserID,
				RSSItemID:      rssItem.ID,
				SubscriptionID: subscription.ID,
			})
		}

		if err != nil {
			logger.Error().Msgf("Failed to write items: %v", err)
		}
	}

	return nil
}
