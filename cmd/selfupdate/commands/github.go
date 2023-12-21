package commands

import (
	"io"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/pkg/cli"
)

func githubCmd() *cli.Command {
	return &cli.Command{
		Name:  "github",
		Usage: "github provider for selfupdate",
		Subcommands: []*cli.Command{
			githubReleaseCmd(),
			githubDownloadCmd(),
		},
	}
}

func githubReleaseCmd() *cli.Command {
	var (
		owner   string
		repor   string
		name    string
		version string
	)

	var sign bool

	return &cli.Command{
		Name:  "release",
		Usage: "release a new version",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "owner",
				Usage:       "owner of the repository",
				Required:    true,
				Destination: &owner,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "name of the repository",
				Required:    true,
				Destination: &repor,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of the binary",
				Required:    true,
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "version of the binary",
				Required:    true,
				Destination: &version,
			},
			&cli.BoolFlag{
				Name:        "sign",
				Usage:       "sign the binary before uploading",
				Destination: &sign,
			},
		},
		Action: func(ctx *cli.Context) error {
			ghClient, err := getGithubClient(owner, repor)
			if err != nil {
				return err
			}

			var r io.Reader = os.Stdin
			if sign {
				privateKey, err := getPrivateKey()
				if err != nil {
					return err
				}

				r = selfupdate.NewHashSigner(privateKey).Sign(ctx.Context, r)
			}

			err = ghClient.Upload(ctx.Context, name, version, r)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func githubDownloadCmd() *cli.Command {
	var (
		owner   string
		repor   string
		name    string
		version string
	)

	var verify bool

	return &cli.Command{
		Name:  "release",
		Usage: "release a new version",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "owner",
				Usage:       "owner of the repository",
				Required:    true,
				Destination: &owner,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "name of the repository",
				Required:    true,
				Destination: &repor,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name of the binary",
				Required:    true,
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "version of the binary",
				Required:    true,
				Destination: &version,
			},
			&cli.BoolFlag{
				Name:        "verify",
				Usage:       "verify the binary after downloading",
				Destination: &verify,
			},
		},
		Action: func(ctx *cli.Context) error {
			ghClient, err := getGithubClient(owner, repor)
			if err != nil {
				return err
			}

			rc := ghClient.Download(ctx.Context, name, version)
			defer rc.Close()

			var r io.Reader = rc

			if verify {
				publicKey, err := getPublicKey()
				if err != nil {
					return err
				}

				r = selfupdate.NewHashVerifier(publicKey).Verify(ctx.Context, r)
				if err != nil {
					return err
				}
			}

			_, err = io.Copy(os.Stdout, r)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func getGithubClient(owner string, repo string) (*selfupdate.Github, error) {
	ghToken, ok := os.LookupEnv("SELF_UPDATE_GH_TOKEN")
	if !ok {
		return nil, cli.Exit("SELF_UPDATE_GH_TOKEN env variable is not set", 1)
	}

	return selfupdate.NewGithub(ghToken, owner, repo), nil
}
