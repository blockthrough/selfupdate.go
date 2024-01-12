package commands

import (
	"fmt"
	"io"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/pkg/cli"
	"selfupdate.blockthrough.com/pkg/crypto"
)

var githubFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "owner",
		Usage:    "owner of the repository",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "repo",
		Usage:    "name of the repository",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "version",
		Usage:    "version of the binary",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "token",
		Usage:    "github repo token, usually provided by github action as GITHUB_TOKEN env",
		Required: true,
	},
}

func githubCmd() *cli.Command {
	return &cli.Command{
		Name:  "github",
		Usage: "a provider tool for working with github api for releasing, uploading and downloading binaries",
		Flags: githubFlags,
		Subcommands: []*cli.Command{
			githubCheckCmd(),
			githubReleaseCmd(),
			githubUploadCmd(),
			githubDownloadCmd(),
		},
	}
}

func githubCheckCmd() *cli.Command {
	var githubCheckFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "filename",
			Usage:    "filename of the binary",
			Required: true,
		},
	}

	return &cli.Command{
		Name:  "check",
		Usage: "check if there is a new version",
		Flags: cli.MergeFlags(githubFlags, githubCheckFlags),
		Action: func(ctx *cli.Context) error {
			owner := ctx.String("owner")
			repo := ctx.String("repo")
			filename := ctx.String("filename")
			version := ctx.String("version")
			ghToken := ctx.String("token")

			ghClient, err := getGithubClient(owner, repo, ghToken)
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
	var githubReleaseFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "title",
			Usage:    "title of the release",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "desc",
			Usage:    "description of the release",
			Required: false,
		},
	}

	return &cli.Command{
		Name:  "release",
		Usage: "create a new github release",
		Flags: cli.MergeFlags(githubFlags, githubReleaseFlags),
		Action: func(ctx *cli.Context) error {
			owner := ctx.String("owner")
			repo := ctx.String("repo")
			version := ctx.String("version")
			ghToken := ctx.String("token")

			title := ctx.String("title")
			desc := ctx.String("desc")

			ghClient, err := getGithubClient(owner, repo, ghToken)
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
	var githubUploadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "filename",
			Usage:    "filename of the binary",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "if provided, it will be used to sign the content before uploading",
		},
	}

	return &cli.Command{
		Name:  "upload",
		Usage: "upload a new asset to an already created github release",
		Flags: cli.MergeFlags(githubFlags, githubUploadFlags),
		Action: func(ctx *cli.Context) error {
			owner := ctx.String("owner")
			repo := ctx.String("repo")
			filename := ctx.String("filename")
			version := ctx.String("version")
			ghToken := ctx.String("token")

			key := ctx.String("key")

			ghClient, err := getGithubClient(owner, repo, ghToken)
			if err != nil {
				return err
			}

			var r io.Reader = os.Stdin
			if key != "" {
				privateKey, err := crypto.ParsePrivateKey(key)
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
	var githubDownloadFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     "filename",
			Usage:    "filename of the binary",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "if provided it will be used to verify the content after downloading",
		},
	}

	return &cli.Command{
		Name:  "download",
		Usage: "download a file from github release's asset",
		Flags: cli.MergeFlags(githubFlags, githubDownloadFlags),
		Action: func(ctx *cli.Context) error {
			owner := ctx.String("owner")
			repo := ctx.String("repo")
			filename := ctx.String("filename")
			version := ctx.String("version")
			ghToken := ctx.String("token")

			key := ctx.String("key")

			ghClient, err := getGithubClient(owner, repo, ghToken)
			if err != nil {
				return err
			}

			rc := ghClient.Download(ctx.Context, filename, version)
			defer rc.Close()

			var r io.Reader = rc

			if key != "" {
				publicKey, err := crypto.ParsePublicKey(key)
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

func getGithubClient(owner string, repo string, token string) (*selfupdate.Github, error) {
	if token == "" {
		return nil, cli.Exit("github token is empty", 1)
	}

	return selfupdate.NewGithub(token, owner, repo), nil
}
