package tasks

import (
	"os"

	"github.com/hibiken/asynq"
)

func Client() *asynq.Client {
	redis := os.Getenv("REDIS_URL")
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redis})

	return asynqClient
}
