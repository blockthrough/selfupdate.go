package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/cmd/selfupdate/commands"
)

// During the build process, these variables are set by Github Actions
// NOTE: if Version is empty or SELF_UPDATE_GH_TOKEN is not set, selfupdating is disabled
var (
	Version   = ""
	PublicKey = ""
)

func main() {
	runUpdate()

	err := commands.Execute(Version)
	if err != nil {
		log.Fatal("Failed to execute: ", err)
		os.Exit(1)
	}
}

func runUpdate() {
	// In order for selfupdating to work, the following conditions must be met:
	// 1. Version must be set
	// 2. SELF_UPDATE_GH_TOKEN must be set
	// 3. PublicKey must be set
	// for setting up the token please refer to
	// "Create a Fine-Grained Personal Access Tokens" in README.md
	ghToken, ok := os.LookupEnv("SELF_UPDATE_GH_TOKEN")
	if !ok {
		fmt.Fprintf(os.Stderr, "Warning: SELF_UPDATE_GH_TOKEN env is not set, selfupdating is disabled\n")
		return
	}

	selfupdate.Auto(
		context.Background(), // Context
		"blockthrough",       // Owner Name
		"selfupdate.go",      // Repo Name
		Version,              // Current Version
		"selfupdate",         // Executable Name,
		ghToken,              // Github Token
		PublicKey,            // Public Key
	)
}
