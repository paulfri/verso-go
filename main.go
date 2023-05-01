package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
	"github.com/versolabs/verso/core/action"
	"github.com/versolabs/verso/server"
	"github.com/versolabs/verso/worker"
)

func main() {
	config := config.GetConfig()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Run the HTTP server",
				Action:  server.Serve(config),
			},
			{
				Name:    "worker",
				Aliases: []string{"w"},
				Usage:   "Run the background worker",
				Action:  worker.Work(config),
			},
			{
				Name:   "crawl",
				Usage:  "Queue the given feed for crawling",
				Action: action.Crawl(config),
			},
			{
				Name:   "config",
				Usage:  "Retrieve configuration",
				Action: action.Config(config),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %+v\n", err)
	}
}
