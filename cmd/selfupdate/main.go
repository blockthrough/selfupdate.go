package main

import (
	"log"
	"os"

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

func main() {
	err := commands.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
