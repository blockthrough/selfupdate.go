package selfupdate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"

	"selfupdate.blockthrough.com/pkg/compress"
)

var (
	ErrGithubReleaseNotFound = errors.New("github release not found")
	ErrGithubRedirect        = errors.New("github redirect")
)

type Github struct {
	owner  string
	name   string
	client *github.Client
}

var _ Uploader = (*Github)(nil)
var _ Checker = (*Github)(nil)
var _ Downloader = (*Github)(nil)

func (g *Github) Upload(ctx context.Context, filename string, version string, r io.Reader) error {
	releaseId, err := g.GetReleaseID(ctx, version)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("repos/%s/%s/releases/%d/assets?name=%s", g.owner, g.name, releaseId, filename)

	req, err := g.client.NewUploadRequest(url, compress.Zip(r), 0, "")
	if err != nil {
		return err
	}

	_, err = g.client.Do(ctx, req, nil)
	return err
}

func (g *Github) Check(ctx context.Context, currentVersion string) (newVersion string, desc string, err error) {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.name, nil)
	if err != nil {
		return
	}

	sort.Slice(releases, func(i, j int) bool {
		return releases[i].GetTagName() > releases[j].GetTagName()
	})

	if len(releases) == 0 || releases[0].GetTagName() == currentVersion {
		return "", "", ErrNoNewVersion
	}

	return releases[0].GetTagName(), releases[0].GetBody(), nil
}

func (g *Github) Download(ctx context.Context, name string, version string) io.ReadCloser {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.name, nil)
	if err != nil {
		return newErrorReader(err)
	}

	var release *github.RepositoryRelease

	for i, _ := range releases {
		if releases[i].GetTagName() == version {
			release = releases[i]
			break
		}
	}

	if release == nil {
		return newErrorReader(ErrGithubReleaseNotFound)
	}

	var githubAsset *github.ReleaseAsset
	for _, asset := range release.Assets {
		if asset.GetName() == name {
			githubAsset = asset
			break
		}
	}

	if githubAsset == nil {
		return newErrorReader(ErrGithubReleaseNotFound)
	}

	rc, redirectURL, err := g.client.Repositories.DownloadReleaseAsset(ctx, g.owner, g.name, githubAsset.GetID(), nil)
	if err != nil {
		return newErrorReader(err)
	}

	// TODO: handle redirect
	// for now we just ignore it
	if redirectURL != "" {
		return newErrorReader(ErrGithubRedirect)
	}

	return compress.Unzip(rc)
}

func (g *Github) GetReleaseID(ctx context.Context, version string) (int64, error) {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.name, nil)
	if err != nil {
		return 0, err
	}

	for _, release := range releases {
		if release.GetTagName() == version {
			return release.GetID(), nil
		}
	}

	return 0, ErrGithubReleaseNotFound
}

func NewGithub(token, repoOwner, repoName string) *Github {
	return &Github{
		owner: repoOwner,
		name:  repoName,
		client: github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{
					AccessToken: token,
				}),
			),
		),
	}
}
