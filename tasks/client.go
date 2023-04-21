package tasks

import (
	"github.com/hibiken/asynq"
)

func Client(redisUrl string) *asynq.Client {
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisUrl})

	return asynqClient
}
