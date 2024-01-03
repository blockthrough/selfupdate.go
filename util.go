package selfupdate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	"selfupdate.blockthrough.com/pkg/crypto"
	"selfupdate.blockthrough.com/pkg/env"
	"selfupdate.blockthrough.com/pkg/executil"
)

func Exec(ctx context.Context, owner string, repo string, currentVersion string, filename string, ext string) {
	executil.Cleanup(ext)

	if currentVersion == "" {
		return
	}

	ghToken, ok := env.Lookup("SELF_UPDATE_GH_TOKEN")
	if !ok {
		fmt.Fprintln(os.Stderr, "SELF_UPDATE_GH_TOKEN env variable is not set")
		return
	}

	publicKeyEnv, ok := env.Lookup("SELF_UPDATE_PUBLIC_KEY")
	if !ok {
		fmt.Fprintln(os.Stderr, "SELF_UPDATE_PUBLIC_KEY env variable is not set")
		return
	}

	publicKey, err := crypto.ParsePublicKey(publicKeyEnv)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to parse public key: %s", err.Error()))
		return
	}

	ghClient := NewGithub(ghToken, owner, repo)

	signedFilename := getSignedFilename(filename)

	newVersion, _, err := ghClient.Check(ctx, signedFilename, currentVersion)
	if errors.Is(err, ErrNoNewVersion) {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Fprintf(os.Stderr, "downloading new version (%s)...", newVersion)

	rc := ghClient.Download(ctx, signedFilename, newVersion)
	defer rc.Close()

	r := NewHashVerifier(publicKey).Verify(ctx, rc)

	err = NewPatcher(ext).Patch(context.Background(), r)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to patch: %s\n", err.Error()))
		return
	}

	fmt.Fprint(os.Stderr, "done\n")

	fmt.Fprintln(os.Stderr, "running new version...")

	err = NewAutoCliRunner(ext).Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getSignedFilename(filename string) string {
	return fmt.Sprintf("%s-%s-%s.sign", filename, runtime.GOOS, runtime.GOARCH)
}
