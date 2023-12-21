package main

import (
	"log"
	"os"

	"selfupdate.blockthrough.com/cmd/selfupdate/commands"
)

// SELF_UPDATE_GH_TOKEN=
// SELF_UPDATE_PRIVATE_KEY=

// SELF_UPDATE_PUBLIC_KEY_FILE_PATH=

// selfupdate generate key
// selfupdate github release --owner blockthrough --repo up-marble --name btctl --version v1.0.0 --sign < ./bin/btctl
// selfupdate github download --owner blockthrough --repo up-marble --name btctl --version v1.0.0 --verify > ./bin/btctl

func main() {
	err := commands.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
