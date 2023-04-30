package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/db/query"
	"github.com/versolabs/verso/util"
	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_PASSWORD = "rectoverso"

func Seed(config *util.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		ctx := context.Background()
		_, queries := Init(config.DatabaseURL, false)

		feed, err1 := queries.CreateRSSFeed(ctx, query.CreateRSSFeedParams{
			Title: "Sounder at Heart",
			URL:   "https://www.sounderatheart.com/rss/current.xml",
		})

		_, err2 := queries.CreateRSSFeed(ctx, query.CreateRSSFeedParams{
			Title: "Sound of Hockey",
			URL:   "https://soundofhockey.com/feed/",
		})

		_, err3 := queries.CreateRSSFeed(ctx, query.CreateRSSFeedParams{
			Title: "Paul's blog",
			URL:   "https://blog.paulfri.xyz/atom.xml",
		})

		password, err4 := bcrypt.GenerateFromPassword([]byte(DEFAULT_PASSWORD), 8)
		user, err5 := queries.CreateUser(ctx, query.CreateUserParams{
			Email:     "test@verso.so",
			Name:      "Verso Test",
			Password:  sql.NullString{String: string(password), Valid: true},
			Superuser: true,
		})

		_, err6 := queries.CreateReaderToken(ctx, query.CreateReaderTokenParams{
			UserID:     user.ID,
			Identifier: "F2vwA2wKSHISLXT7slqt",
		})

		_, err7 := queries.CreateRSSSubscription(
			ctx,
			query.CreateRSSSubscriptionParams{
				UserID: user.ID,
				FeedID: feed.ID,
			},
		)

		tag, err8 := queries.CreateTag(
			ctx,
			query.CreateTagParams{
				UserID: user.ID,
				Name:   "soccer",
			},
		)

		_, err9 := queries.CreateRSSFeedTag(
			ctx,
			query.CreateRSSFeedTagParams{
				TagID:     tag.ID,
				RSSFeedID: feed.ID,
			},
		)

		err := errors.Join(err1, err2, err3, err4, err5, err6, err7, err8, err9)

		if err != nil {
			log.Fatalf("failed to seed database: %v", err)
		}

		return nil
	}
}
