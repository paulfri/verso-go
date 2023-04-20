package cli

import (
	"strconv"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/tasks"
)

func Crawl(cliContext *cli.Context) error {
	client := tasks.Client()
	feedId := cliContext.Args().Get(0)

	i, err := strconv.ParseInt(feedId, 10, 64)
	if err != nil {
		panic(err)
	}

	task, _ := tasks.NewFeedParseTask(i)
	_, err2 := client.Enqueue(task)

	return err2
}
