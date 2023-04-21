package tasks

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/util"
)

const (
	TypeFeedParse = "feed:parse"
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
		worker := Worker{
			Container: &util.Container{
				Asynq:   Client(config.RedisUrl),
				DB:      database,
				Queries: queries,
				Render:  render.New(),
			},
		}

		mux := asynq.NewServeMux()
		mux.HandleFunc(TypeFeedParse, worker.HandleFeedParseTask)

		if err := srv.Run(mux); err != nil {
			log.Fatal(err)

			return err
		}

		return nil
	}
}
