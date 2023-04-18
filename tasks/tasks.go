package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeFeedParse = "feed:parse"
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

	return nil
}
