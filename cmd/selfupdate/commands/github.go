package commands

import "selfupdate.blockthrough.com/pkg/cli"

func githubCmd() *cli.Command {
	return &cli.Command{
		Name:        "github",
		Usage:       "github provider for selfupdate",
		Subcommands: []*cli.Command{},
	}
}
