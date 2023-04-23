package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker/tasks"
)

type Worker struct {
	Container *util.Container
}

func Work(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.RedisURL},
			asynq.Config{Concurrency: config.WorkerConcurrency},
		)

		database, queries := db.Init(config.DatabaseURL)
		client := Client(config.RedisURL)

		worker := Worker{
			Container: &util.Container{
				Asynq:   client,
				DB:      database,
				Queries: queries,
				Render:  render.New(),
			},
		}

		mux := asynq.NewServeMux()
		mux.HandleFunc(tasks.TypeFeedParse, worker.HandleFeedParseTask)

		// TODO: enable
		// worker.ConfigureCron()

		if err := srv.Run(mux); err != nil {
			log.Fatal(err)
		}

		return nil
	}
}
