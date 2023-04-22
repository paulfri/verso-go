package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/versolabs/citra/db/query"
)

type FeedParsePayload struct {
	FeedId int64
}

func NewFeedParseTask(feedId int64) (*asynq.Task, error) {
	payload, err := json.Marshal(FeedParsePayload{FeedId: feedId})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeFeedParse, payload), nil
}

func (worker Worker) HandleFeedParseTask(ctx context.Context, t *asynq.Task) error {
	var p FeedParsePayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Parsing feed: feed_id=%d\n", p.FeedId)

	thisFeed, err := worker.Container.Queries.FindRssFeed(ctx, int64(p.FeedId))
	if err != nil {
		return err
	}

	fmt.Println(thisFeed)

	url := thisFeed.Url

	parser := gofeed.NewParser()
	feed, _ := parser.ParseURL(url)
	for _, item := range feed.Items {
		fmt.Println(item.Title)

		_, err := worker.Container.Queries.CreateRssItem(ctx, query.CreateRssItemParams{
			RssFeedID:   int64(p.FeedId),
			RssGuid:     item.GUID,
			Title:       item.Title,
			Content:     item.Content,
			Link:        item.Link,
			PublishedAt: sql.NullTime{Time: *item.PublishedParsed, Valid: true},
			// RemoteUpdatedAt: sql.NullTime{Time: *item.UpdatedParsed, Valid: true},
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
