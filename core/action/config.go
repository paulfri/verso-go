package action

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/versolabs/verso/config"
)

func Config(config *config.Config) cli.ActionFunc {
	return func(cliContext *cli.Context) error {
		switch cliContext.Args().Get(0) {
		case "dbname":
			fmt.Println(config.Database.Database)
		case "dbconn":
			fmt.Println(config.Database.Conn)
		case "dburl":
			fmt.Println(config.Database.URL)
		default:
			panic("Unknown argument")
		}

		return nil
	}
}
