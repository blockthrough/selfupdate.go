package selfupdate

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNoNewVersion = errors.New("no new version")
)

type Signer interface {
	Sign(ctx context.Context, r io.Reader) io.Reader
}

type Uploader interface {
	Upload(ctx context.Context, filename string, version string, r io.Reader) error
}

type Checker interface {
	Check(ctx context.Context, filename string, currentVersion string) (newVersion string, desc string, err error)
}

type Downloader interface {
	Download(ctx context.Context, name string, version string) io.ReadCloser
}

type Verifier interface {
	Verify(ctx context.Context, r io.Reader) io.Reader
}

type Patcher interface {
	Patch(ctx context.Context, patch io.Reader) error
}

type Runner interface {
	Run(ctx context.Context) error
}

// helper types

type SignerFunc func(ctx context.Context, r io.Reader) io.Reader

var _ Signer = SignerFunc(nil)

func (f SignerFunc) Sign(ctx context.Context, r io.Reader) io.Reader {
	return f(ctx, r)
}

type VerifierFunc func(ctx context.Context, r io.Reader) io.Reader

var _ Verifier = VerifierFunc(nil)

func (f VerifierFunc) Verify(ctx context.Context, r io.Reader) io.Reader {
	return f(ctx, r)
}

type RunnerFunc func(ctx context.Context) error

var _ Runner = RunnerFunc(nil)

func (f RunnerFunc) Run(ctx context.Context) error {
	return f(ctx)
}

type PatcherFunc func(ctx context.Context, patch io.Reader) error

var _ Patcher = PatcherFunc(nil)

func (f PatcherFunc) Patch(ctx context.Context, patch io.Reader) error {
	return f(ctx, patch)
}
