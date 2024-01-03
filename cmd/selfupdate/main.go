package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/cmd/selfupdate/commands"
	"selfupdate.blockthrough.com/pkg/crypto"
	"selfupdate.blockthrough.com/pkg/env"
)

// SELF_UPDATE_GH_TOKEN=
// SELF_UPDATE_PRIVATE_KEY=
// SELF_UPDATE_PUBLIC_KEY=

// selfupdate crypto generate-keys
// selfupdate crypto sign < ./bin/btctl > ./bin/btctl.sig
// selfupdate crypto verify < ./bin/btctl.sig > ./bin/btctl
// selfupdate github release --owner blockthrough --repo up-marble --name btctl.sign --version v1.0.0 --sign < ./bin/btctl
// selfupdate github download --owner blockthrough --repo up-marble --name btctl.sign --version v1.0.0 --verify > ./bin/btctl

var Version string = ""
var OS string = ""
var Arch string = ""

func main() {
	updateAndRun(context.Background())

	err := commands.Execute(Version)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func updateAndRun(ctx context.Context) {
	if Version == "" {
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

	ghClient := selfupdate.NewGithub(ghToken, "blockthrough", "selfupdate.go")

	newVersion, _, err := ghClient.Check(ctx, getSignFilename(), Version)
	if errors.Is(err, selfupdate.ErrNoNewVersion) {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Fprintf(os.Stderr, "downloading new version (%s)...", newVersion)

	rc := ghClient.Download(ctx, getSignFilename(), newVersion)
	defer rc.Close()

	r := selfupdate.NewHashVerifier(publicKey).Verify(ctx, rc)

	err = selfupdate.NewPatcher().Patch(context.Background(), r)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("failed to patch: %s\n", err.Error()))
		return
	}

	fmt.Fprint(os.Stderr, "done\n")

	fmt.Fprintln(os.Stderr, "running new version...")

	err = selfupdate.NewAutoCliRunner().Run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getSignFilename() string {
	return fmt.Sprintf("selfupdate-%s-%s.sign", OS, Arch)
}
