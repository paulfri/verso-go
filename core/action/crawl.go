package action

import (
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
	"github.com/versolabs/verso/worker"
	"github.com/versolabs/verso/worker/tasks"
)

func Crawl(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		client := worker.Client(config.RedisURL)
		feedID := cliContext.Args().Get(0)

		i, err := strconv.ParseInt(feedID, 10, 64)
		if err != nil {
			panic(err)
		}

		task, _ := tasks.NewFeedParseTask(i)
		_, err2 := client.Enqueue(task)

		return err2
	}
}
