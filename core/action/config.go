package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
)

func Config(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		var sslMode string
		if config.Database.SSLDisabled {
			sslMode = "disable"
		} else {
			sslMode = "require"
		}

		switch cliContext.Args().Get(0) {
		case "dbname":
			fmt.Println(config.Database.Database)
		case "database":
			fmt.Printf(
				"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
				config.Database.User,
				config.Database.Password,
				config.Database.Host,
				config.Database.Port,
				config.Database.Database,
				sslMode,
			)
		case "goose":
			fmt.Printf(
				"host=%s user=%s password=%s dbname=%s sslmode=%s",
				config.Database.Host,
				config.Database.User,
				config.Database.Password,
				config.Database.Database,
				sslMode,
			)
		default:
			panic("Unknown argument")
		}

		return nil
	}
}
