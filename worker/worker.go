package worker

import (
	"context"
	"log"

	"github.com/airbrake/gobrake/v5"
	"github.com/hibiken/asynq"
	"github.com/unrolled/render"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker/tasks"
)

type Worker struct {
	Container *util.Container
}

func Work(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		airbrake := util.Airbrake(config)
		notifier := notify(airbrake)
		handler := asynq.ErrorHandlerFunc(notifier)

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: config.RedisURL},
			asynq.Config{
				Concurrency:  config.Worker.Concurrency,
				ErrorHandler: handler,
			},
		)

		database, queries := db.Init(config.Database.URL(), false)
		client := Client(config.RedisURL)

		worker := Worker{
			Container: &util.Container{

				Asynq:   client,
				DB:      database,
				Queries: queries,
				Render:  render.New(),
				Logger:  util.Logger(),
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

func notify(airbrake *gobrake.Notifier) asynq.ErrorHandlerFunc {
	return func(ctx context.Context, task *asynq.Task, err error) {
		airbrake.Notify(err, nil)
	}
}
