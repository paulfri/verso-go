package action

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
)

func Config(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		switch cliContext.Args().Get(0) {
		case "dbname":
			os.Stdout.Write([]byte(config.Database.Database))
		case "dbconn":
			os.Stdout.Write([]byte(config.Database.Conn))
		case "dburl":
			os.Stdout.Write([]byte(config.Database.URL))
		default:
			panic("Unknown argument")
		}

		return nil
	}
}
