package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db/query"
)

func Seed(cliContext *cli.Context) error {
	context := context.Background()
	queries := Queries()

	_, err1 := queries.CreateRssFeed(context, query.CreateRssFeedParams{
		Title: "Sounder at Heart",
		Url:   "https://www.sounderatheart.com/rss/current.xml",
	})

	_, err2 := queries.CreateRssFeed(context, query.CreateRssFeedParams{
		Title: "Sound of Hockey",
		Url:   "https://soundofhockey.com/feed/",
	})

	err := errors.Join(err1, err2)

	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	return nil
}
