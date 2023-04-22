package cli

import (
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/util"
	"github.com/versolabs/citra/worker"
	"github.com/versolabs/citra/worker/tasks"
)

func Crawl(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		client := worker.Client(config.RedisUrl)
		feedId := cliContext.Args().Get(0)

		i, err := strconv.ParseInt(feedId, 10, 64)
		if err != nil {
			panic(err)
		}

		task, _ := tasks.NewFeedParseTask(i)
		_, err2 := client.Enqueue(task)

		return err2
	}
}
