package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeFeedParse = "feed:parse"
)

type FeedParsePayload struct {
	FeedID int64
}

func NewFeedParseTask(feedID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(FeedParsePayload{FeedID: feedID})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeFeedParse, payload), nil
}
