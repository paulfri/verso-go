package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/ai"
	"github.com/versolabs/citra/feed"
)

func main() {
	app := &cli.App{
		Name:  "citra",
		Usage: "Summarize the latest item from a given RSS feed",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "feed",
				Aliases: []string{"f"},
				// Value: "https://www.sounderatheart.com/rss/current.xml",
				Usage:    "Feed to fetch and summarize",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			items := feed.Fetch(ctx.String("feed"))
			summary := ai.Summarize(items[0].Content)

			fmt.Println(summary)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
