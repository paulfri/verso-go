package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/versolabs/citra/db"
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
	thisFeed, err := queries.GetFeedById(ctx, int32(p.FeedId))
	if err != nil {
		return err
	}

	fmt.Println(thisFeed)

	url := thisFeed.Url

	items := feed.Fetch(url)
	fmt.Println(lo.Map(items, func(item *gofeed.Item, index int) string {
		return item.Title + "\n"
	}))

	return nil
}
