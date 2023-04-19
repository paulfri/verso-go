package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/feed"
)

type FeedParsePayload struct {
	FeedId int
}

func NewFeedParseTask(feedId int) (*asynq.Task, error) {
	payload, err := json.Marshal(FeedParsePayload{FeedId: feedId})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeFeedParse, payload), nil
}

func HandleFeedParseTask(ctx context.Context, t *asynq.Task) error {
	var p FeedParsePayload

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	fmt.Printf("Parsing feed: feed_id=%d\n", p.FeedId)

	queries := db.Queries()
	thisFeed, err := queries.GetRssFeedById(ctx, int64(p.FeedId))
	if err != nil {
		return err
	}

	fmt.Println(thisFeed)

	url := thisFeed.Url

	items := feed.Fetch(url)
	for _, item := range items {
		fmt.Println(item.Title)

		_, err := queries.CreateItem(ctx, query.CreateItemParams{
			RssFeedID:       int64(p.FeedId),
			RssGuid:         item.GUID,
			Title:           item.Title,
			Content:         item.Content,
			Link:            item.Link,
			PublishedAt:     sql.NullTime{Time: *item.PublishedParsed, Valid: true},
			RemoteUpdatedAt: sql.NullTime{Time: *item.UpdatedParsed, Valid: true},
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
