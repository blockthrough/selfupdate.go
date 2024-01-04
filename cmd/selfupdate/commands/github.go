package commands

import (
	"fmt"
	"io"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/pkg/cli"
	"selfupdate.blockthrough.com/pkg/env"
)

func githubCmd() *cli.Command {
	return &cli.Command{
		Name:  "github",
		Usage: "a provider tool for working with github api for releasing, uploading and downloading binaries",
		Subcommands: []*cli.Command{
			githubCheckCmd(),
			githubReleaseCmd(),
			githubUploadCmd(),
			githubDownloadCmd(),
		},
	}
}

func githubCheckCmd() *cli.Command {
	var (
		owner    string
		repo     string
		filename string
		version  string
	)

	return &cli.Command{
		Name:  "check",
		Usage: "check if there is a new version",
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
				Destination: &repo,
			},
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "name of the binary file",
				Required:    true,
				Destination: &filename,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "version of the binary",
				Required:    true,
				Destination: &version,
			},
		},
		Action: func(ctx *cli.Context) error {
			ghClient, err := getGithubClient(owner, repo)
			if err != nil {
				return err
			}

			newVersion, desc, err := ghClient.Check(ctx.Context, filename, version)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "new version: %s\n", newVersion)
			fmt.Fprintf(os.Stdout, "description: %s\n", desc)

			return nil
		},
	}
}

func githubReleaseCmd() *cli.Command {
	var (
		owner   string
		repo    string
		version string
		title   string
		desc    string
	)

	return &cli.Command{
		Name:  "release",
		Usage: "create a new github release",
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
				Destination: &repo,
			},
			&cli.StringFlag{
				Name:        "version",
				Usage:       "version of the binary",
				Required:    true,
				Destination: &version,
			},
			&cli.StringFlag{
				Name:        "title",
				Usage:       "title of the release",
				Required:    true,
				Destination: &title,
			},
			&cli.StringFlag{
				Name:        "desc",
				Usage:       "description of the release",
				Required:    false,
				Destination: &desc,
			},
		},
		Action: func(ctx *cli.Context) error {
			ghClient, err := getGithubClient(owner, repo)
			if err != nil {
				return err
			}

			// NOTE: this check makes sure we are not creating a release that already exists
			// which leads to error and make the cli return and error. This is usually not a
			// problem unless the cli gets executed in github actions with strategy matrix.
			exists, err := ghClient.CheckIfReleaseExists(ctx.Context, version)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("release %s already exists", version)
			}

			err = ghClient.Release(ctx.Context, version, title, desc)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func githubUploadCmd() *cli.Command {
	var (
		owner    string
		repo     string
		filename string
		version  string
	)

	var sign bool

	return &cli.Command{
		Name:  "upload",
		Usage: "upload a new asset to an already created github release",
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
				Destination: &repo,
			},
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "filename of the binary",
				Required:    true,
				Destination: &filename,
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
			ghClient, err := getGithubClient(owner, repo)
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

			err = ghClient.Upload(ctx.Context, filename, version, r)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func githubDownloadCmd() *cli.Command {
	var (
		owner    string
		repo     string
		filename string
		version  string
	)

	var verify bool

	return &cli.Command{
		Name:  "download",
		Usage: "download a file from github release's asset",
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
				Destination: &repo,
			},
			&cli.StringFlag{
				Name:        "filename",
				Usage:       "name of the binary content",
				Required:    true,
				Destination: &filename,
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
			ghClient, err := getGithubClient(owner, repo)
			if err != nil {
				return err
			}

			rc := ghClient.Download(ctx.Context, filename, version)
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
	ghToken, ok := env.Lookup("SELF_UPDATE_GH_TOKEN")
	if !ok {
		return nil, cli.Exit("SELF_UPDATE_GH_TOKEN env variable is not set", 1)
	}

	return selfupdate.NewGithub(ghToken, owner, repo), nil
}
