package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	citraCli "github.com/versolabs/citra/cli"
	"github.com/versolabs/citra/db"
	"github.com/versolabs/citra/http"
	"github.com/versolabs/citra/tasks"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Run the HTTP server",
				Action:  http.Serve,
			},
			{
				Name:    "worker",
				Aliases: []string{"s"},
				Usage:   "Run the background worker",
				Action:  tasks.Work,
			},
			{
				Name:   "seed",
				Usage:  "Seed the database with test fixtures",
				Action: db.Seed,
			},
			{
				Name:   "crawl",
				Usage:  "Queue the given feed for crawling",
				Action: citraCli.Crawl,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
