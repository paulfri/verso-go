package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
)

func Config(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		switch cliContext.Args().Get(0) {
		case "database":
			fmt.Printf("%s\n", config.Database.Database)
		case "goose":
			fmt.Printf(
				"host=%s user=%s password=%s dbname=%s sslmode=disable",
				config.Database.Host,
				config.Database.User,
				config.Database.Password,
				config.Database.Database,
			)
		default:
			panic("Unknown argument")
		}

		return nil
	}
}
