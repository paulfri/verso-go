package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/citra/db/query"
	"github.com/versolabs/citra/util"
	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_PASSWORD = "rectoverso"

func Seed(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		context := context.Background()
		_, queries := Init(config.DatabaseURL)

		_, err1 := queries.CreateRSSFeed(context, query.CreateRSSFeedParams{
			Title: "Sounder at Heart",
			URL:   "https://www.sounderatheart.com/rss/current.xml",
		})

		_, err2 := queries.CreateRSSFeed(context, query.CreateRSSFeedParams{
			Title: "Sound of Hockey",
			URL:   "https://soundofhockey.com/feed/",
		})

		password, err3 := bcrypt.GenerateFromPassword([]byte(DEFAULT_PASSWORD), 8)
		_, err4 := queries.CreateUser(context, query.CreateUserParams{
			Email:     "paul@verso.so",
			Name:      "Paul Friedman",
			Password:  sql.NullString{String: string(password), Valid: true},
			Superuser: true,
		})

		err := errors.Join(err1, err2, err3, err4)

		if err != nil {
			fmt.Println(err)
			panic(1)
		}

		return nil
	}
}
