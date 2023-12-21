package commands

import (
	"os"

	"selfupdate.blockthrough.com/pkg/cli"
)

func Execute() error {
	app := cli.App{
		Name:    "selfupdate",
		Usage:   "a cli for self-update of golang apps",
		Version: "1.0.0",
		Commands: []*cli.Command{
			cryptoCmd(),
			githubCmd(),
		},
	}

	return app.Run(os.Args)
}
