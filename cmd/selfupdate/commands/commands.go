package commands

import (
	"os"

	"selfupdate.blockthrough.com/pkg/cli"
)

func Execute(version string) error {
	app := cli.App{
		Name:    "selfupdate",
		Usage:   "a cli for self-update of golang apps",
		Version: version,
		Commands: []*cli.Command{
			cryptoCmd(),
			githubCmd(),
		},
	}

	return app.Run(os.Args)
}
