package commands

import (
	"os"

	"selfupdate.blockthrough.com/pkg/cli"
)

func Execute() error {
	app := cli.App{
		Name:    "selfupdate",
		Usage:   "golang selfupdate cli",
		Version: "1.0.0",
		Commands: []*cli.Command{
			generateCmd(),
		},
	}

	return app.Run(os.Args)
}
