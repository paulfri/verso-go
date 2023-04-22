package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/util"
	"github.com/versolabs/citra/worker/tasks"
)

type Worker struct {
	Container *util.Container
}

func Work(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.RedisUrl},
			asynq.Config{Concurrency: config.WorkerConcurrency},
		)

		database, queries := db.Init(config.DatabaseUrl)
		worker := Worker{
			Container: &util.Container{
				Asynq:   Client(config.RedisUrl),
				DB:      database,
				Queries: queries,
				Render:  render.New(),
			},
		}

		mux := asynq.NewServeMux()
		mux.HandleFunc(tasks.TypeFeedParse, worker.HandleFeedParseTask)

		if err := srv.Run(mux); err != nil {
			log.Fatal(err)
		}

		return nil
	}
}

func (worker *Worker) HandleFeedParseTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.FeedParsePayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Parsing feed: feed_id=%d\n", p.FeedId)

	thisFeed, err := worker.Container.Queries.FindRssFeed(ctx, int64(p.FeedId))
	if err != nil {
		return err
	}

	fmt.Println(thisFeed)

	// Fetch the feed and parse it.
	url := thisFeed.Url
	parser := gofeed.NewParser()
	remoteFeed, _ := parser.ParseURL(url)

	items := remoteFeed.Items
	if len(items) == 0 {
		// Empty feed. Suspicious.
		// TODO: probably log this error somewhere
		return nil
	}

	subscribers, err := worker.Container.Queries.GetSubscribersByFeedId(ctx, thisFeed.ID)
	if err != nil {
		return err
	}

	for _, feedItem := range remoteFeed.Items {
		fmt.Println(feedItem.Title)

		rssItem, err := worker.Container.Queries.CreateRssItem(ctx, query.CreateRssItemParams{
			RssFeedID:   int64(p.FeedId),
			RssGuid:     feedItem.GUID,
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
				RssItemID: sql.NullInt64{Int64: rssItem.ID, Valid: true},
			})
		}

		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
