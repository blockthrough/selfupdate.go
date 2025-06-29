package selfupdate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"

	"selfupdate.blockthrough.com/pkg/compress"
	"selfupdate.blockthrough.com/pkg/version"
)

var (
	ErrGithubAssetNotFound   = errors.New("github asset not found")
	ErrGithubReleaseNotFound = errors.New("github release not found")
	ErrGithubRedirect        = errors.New("github redirect")
)

type Github struct {
	owner            string
	repo             string
	client           *github.Client
	versionCompareFn func(a, b string) bool
}

var _ Uploader = (*Github)(nil)
var _ Checker = (*Github)(nil)
var _ Downloader = (*Github)(nil)

func (g *Github) Upload(ctx context.Context, filename string, version string, r io.Reader) error {
	// to overrride existing asset, we need to delete it first
	// This is a way if we needed to rerun the github action again
	err := g.DeleteAsset(ctx, filename, version)
	if err != nil {
		return err
	}

	releaseId, err := g.GetReleaseIDByVersion(ctx, version)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("repos/%s/%s/releases/%d/assets?name=%s", g.owner, g.repo, releaseId, filename)

	var buffer bytes.Buffer
	n, err := io.Copy(&buffer, compress.Zip(r))
	if err != nil {
		return err
	}

	req, err := g.client.NewUploadRequest(url, &buffer, n, "")
	if err != nil {
		return err
	}

	_, err = g.client.Do(ctx, req, nil)
	return err
}

func (g *Github) Release(ctx context.Context, tag string, releaseTitle string, releaseBody string) error {
	// 781b176f2d5a4d1887ba386fed2bae0f6ab3bb92
	_, _, err := g.client.Repositories.CreateRelease(ctx, g.owner, g.repo, &github.RepositoryRelease{
		TagName: &tag,
		// TargetCommitish: &targetCommitish,
		Name:       &releaseTitle,
		Body:       &releaseBody,
		Draft:      github.Bool(false),
		Prerelease: github.Bool(false),
	})
	return err
}

func (g *Github) Check(ctx context.Context, filename string, currentVersion string) (newVersion string, desc string, err error) {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.repo, nil)
	if err != nil {
		return
	}

	sort.Slice(releases, func(i, j int) bool {
		return g.versionCompareFn(releases[i].GetTagName(), releases[j].GetTagName())
	})

	if len(releases) == 0 || releases[0].GetTagName() == currentVersion {
		return "", "", ErrNoNewVersion
	}

	release := releases[0]
	var githubAsset *github.ReleaseAsset
	for _, asset := range release.Assets {
		if asset.GetName() == filename {
			githubAsset = asset
			break
		}
	}

	if githubAsset == nil {
		return "", "", ErrGithubAssetNotFound
	}

	return releases[0].GetTagName(), releases[0].GetBody(), nil
}

func (g *Github) DeleteAsset(ctx context.Context, filename string, version string) (err error) {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.repo, nil)
	if err != nil {
		return
	}

	var targetRelease *github.RepositoryRelease
	for _, release := range releases {
		if release.GetTagName() == version {
			targetRelease = release
			break
		}
	}

	if targetRelease == nil {
		return nil
	}

	var existingAssetID int64
	for _, asset := range targetRelease.Assets {
		if asset.GetName() == filename {
			existingAssetID = asset.GetID()
			break
		}
	}

	if existingAssetID == 0 {
		return nil
	}

	_, err = g.client.Repositories.DeleteReleaseAsset(ctx, g.owner, g.repo, existingAssetID)
	if err != nil {
		return err
	}

	return nil
}

func (g *Github) Download(ctx context.Context, name string, version string) io.ReadCloser {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.repo, nil)
	if err != nil {
		return newErrorReader(err)
	}

	var release *github.RepositoryRelease

	for i := range releases {
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

	rc, redirectURL, err := g.client.Repositories.DownloadReleaseAsset(ctx, g.owner, g.repo, githubAsset.GetID(), nil)
	if err != nil {
		return newErrorReader(err)
	}

	if redirectURL == "" {
		return compress.Unzip(rc)
	}

	resp, err := http.Get(redirectURL)
	if err != nil {
		return newErrorReader(err)
	}

	return compress.Unzip(resp.Body)
}

func (g *Github) GetReleaseIDByVersion(ctx context.Context, version string) (int64, error) {
	releases, _, err := g.client.Repositories.ListReleases(ctx, g.owner, g.repo, nil)
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

type githubOptFn func(g *Github)

func WithGithubVersionCompare(fn func(a, b string) bool) githubOptFn {
	return func(g *Github) {
		g.versionCompareFn = fn
	}
}

func NewGithub(token, repoOwner, repoName string, optFns ...githubOptFn) *Github {
	g := &Github{
		owner: repoOwner,
		repo:  repoName,
		client: github.NewClient(
			oauth2.NewClient(
				context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{
					AccessToken: token,
				}),
			),
		),
		versionCompareFn: version.Compare,
	}

	for _, optFn := range optFns {
		optFn(g)
	}

	return g
}
