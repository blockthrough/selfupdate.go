package selfupdate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"selfupdate.blockthrough.com/pkg/crypto"
	"selfupdate.blockthrough.com/pkg/env"
	"selfupdate.blockthrough.com/pkg/executil"
)

func Auto(ctx context.Context, owner string, repo string, currentVersion string, filename string) {
	if currentVersion == "" {
		return
	}

	currentExecPath, err := executil.CurrentPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to get current executable path: %s", err.Error()))
		return
	}

	actualFilename := filepath.Base(currentExecPath)
	actualFileExt := filepath.Ext(actualFilename)
	newFilename := filepath.Join(filepath.Dir(currentExecPath), filename+"-downloaded"+actualFileExt)

	// if the filename is not the same as the current executable, then we are
	// running the patcher executable, so we simply
	if actualFilename != filename {
		// this is a good chance to copy the downloaded file to the original file
		err = copyFile(filepath.Join(currentExecPath, filepath.Dir(currentExecPath), filename+actualFileExt), currentExecPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to copy the downloaded file over original onex: %s", err.Error()))
		}
		return
	} else {
		// this is the actual executable, this is a good time to remove the
		// downloaded file, if there is one, that's why we don't care about
		// the error
		os.Remove(newFilename)
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

	signedFilename := fmt.Sprintf("%s-%s-%s.sign", filename, runtime.GOOS, runtime.GOARCH)

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

	err = NewPatcher(newFilename).Patch(context.Background(), NewHashVerifier(publicKey).Verify(ctx, rc))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to patch: %s\n", err.Error()))
		return
	}

	fmt.Fprint(os.Stderr, "done\n")

	fmt.Fprintln(os.Stderr, "running new version...")

	err = NewCliRunner(newFilename, os.Args[1:]...).Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(0)
}
