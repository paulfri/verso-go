package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	citraCli "github.com/versolabs/citra/cli"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/server"
	"github.com/versolabs/citra/tasks"
	"github.com/versolabs/citra/util"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %+v\n", err)
	}

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
				Aliases: []string{"s"},
				Usage:   "Run the background worker",
				Action:  tasks.Work(&config),
			},
			{
				Name:   "seed",
				Usage:  "Seed the database with test fixtures",
				Action: db.Seed(&config),
			},
			{
				Name:   "crawl",
				Usage:  "Queue the given feed for crawling",
				Action: citraCli.Crawl(&config),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running app: %+v\n", err)
	}
}
