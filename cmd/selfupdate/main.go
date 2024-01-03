package main

import (
	"context"
	"log"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/cmd/selfupdate/commands"
)

// SELF_UPDATE_GH_TOKEN=
// SELF_UPDATE_PRIVATE_KEY=
// SELF_UPDATE_PUBLIC_KEY=

// selfupdate crypto generate-keys
// selfupdate crypto sign < ./bin/btctl > ./bin/btctl.sig
// selfupdate crypto verify < ./bin/btctl.sig > ./bin/btctl
// selfupdate github release --owner blockthrough --repo up-marble --name btctl.sign --version v1.0.0 --sign < ./bin/btctl
// selfupdate github download --owner blockthrough --repo up-marble --name btctl.sign --version v1.0.0 --verify > ./bin/btctl

var Version string = "v0.0.0"

func main() {
	selfupdate.Exec(
		context.Background(), // Context
		"blockthrough",       // Owner Name
		"selfupdate.go",      // Repo Name
		Version,              // Current Version
		"selfupdate",         // Executable Name
		".new",               // Executable Extension for downloading new version
	)

	err := commands.Execute(Version)
	if err != nil {
		log.Fatal("Failed to execute: ", err)
		os.Exit(1)
	}
}
