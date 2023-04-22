package worker

import (
	"github.com/hibiken/asynq"
)

func Client(url string) *asynq.Client {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: url})

	return asynqClient
}
