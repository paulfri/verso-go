package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/util"
	"github.com/versolabs/citra/worker/tasks"
)

type Worker struct {
	Container *util.Container
}

func Work(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.RedisUrl},
			asynq.Config{Concurrency: config.WorkerConcurrency},
		)

		database, queries := db.Init(config.DatabaseUrl)
		client := Client(config.RedisUrl)

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
