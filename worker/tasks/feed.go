package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeFeedParse = "feed:parse"
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
