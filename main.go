package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/core/action"
	"github.com/versolabs/verso/db"
	"github.com/versolabs/verso/server"
	"github.com/versolabs/verso/util"
	"github.com/versolabs/verso/worker"
)

func main() {
	// Load environment from .env file, if present. Production environments are
	// configured without the use of this file.
	_ = godotenv.Load()

	config := util.GetConfig()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Run the HTTP server",
				Action:  server.Serve(&config),
			},
			{
				Name:    "worker",
				Aliases: []string{"w"},
				Usage:   "Run the background worker",
				Action:  worker.Work(&config),
			},
			{
				Name:   "seed",
				Usage:  "Seed the database with test fixtures",
				Action: db.Seed(&config),
			},
			{
				Name:   "crawl",
				Usage:  "Queue the given feed for crawling",
				Action: action.Crawl(&config),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %+v\n", err)
	}
}
