package tasks

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/urfave/cli/v2"
)

func Work(cliContext *cli.Context) error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeFeedParse, HandleFeedParseTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
