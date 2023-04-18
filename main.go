package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
