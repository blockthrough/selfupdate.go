# Self Update for Golang

It's a dead-simple toolchain for updating Golang binaries. It is designed to leverage Github's Releases and Actions. It comes with its CLI and SDK to make things even simpler.

## CLI Installation

```bash
go install selfupdate.blockthough.com@latest
```

## CLI Usage

`selfupdate`'s cli consists of a couple of handy commands which are described as follows

> NOTE: all the helper commands can be accessed via `--help` or `-h`

### crypto

it's a set of tools and commands that help generate public/private keys, sign and verify content

#### keys

generates a pair of public and private keys

```bash
selfupdate crypto keys
```

#### sign

sign a binary using a private key. The private key must be set using the environment variable `SELF_UPDATE_PRIVATE_KEY`

```bash
export SELF_UPDATE_PRIVATE_KEY=....
selfupdate crypto sign < ./bin/file > ./bin/file.sig
```

#### verify

verify a binary using a public key. The public key must be set using the environment variable `SELF_UPDATE_PUBLIC_KEY`

```bash
export SELF_UPDATE_PUBLIC_KEY=....
selfupdate crypto verify < ./bin/file.sg > ./bin/file
```

### github

a provider tool for working with github's apis for releasing, uploading and downloading binaries.

> NOTE: an environment variable, `SELF_UPDATE_GH_TOKEN`, needs to be created and set with `GITHUB_TOKEN`. It is highly recommened to use it `GITHUB_TOKEN` rather than personal token. Also, this command with all its sub commands needs to be called inside github's actions workflow.

#### check

check if there is a new version by providing the version

```bash
selfupdate github check --owner blockthough --repo selfupdate.go --filename selfupdate --version v0.0.1
```

> for more info, run `selfupdate github check --help`

#### release

To provide a better developer experience, this command can create a new github release, with a title, and a description and attach it to a specific tag.

> NOTE: running this command with the same version will be noop. So it can be safely used in github actions workflow with matrix strategy.

```bash
selfupdate github release -owner blockthough --repo selfupdate.go --version v0.0.1 --title version v0.0.1 --desc "this is an amazin release"
```

> for more info, run `selfupdate github release --help`

#### upload

upload a new asset to an already created github release. It is a requirement to set an environment variable, `SELF_UPDATE_PRIVATE_KEY`, while using `--sign` flag.

the `--sign` flag will sign the uploaded content which later can be verified using `SELF_UPDATE_PUBLIC_KEY` environment variable.

> In order to upload assets, a github release must be created first. Please refer to `release` subcommand. Also this command can be used multiple times for each individual asset in github actions workflow.

```bash
selfupdate github upload -owner blockthough --repo selfupdate.go --version v0.0.1 --filename selfupload.sign --sign < /path/to/file
```

#### download

In order to download a specific asset from a gethub release, this command can be used. It requires `--filename` and `--version` to be presented.

```bash
selfupdate github download -owner blockthough --repo selfupdate.go --version v0.0.1 --filename selfupload.sign --verify > /path/to/file
```

## Usage

In order to have successful self-updating binaries, two steps need to be followed:

### ( 1 ) Github Actions Workflow

- First compile the your code and generate a binary, make sure to use `-ldflags "-X main.Version=${{ github.ref_name }}"` flag during `go build` to inject the new tag as a version into the binary.

- Create a new Release using `selfupdate github release` command
- Sign and upload content using `selfupdate github upload`

### ( 2 ) Using `selfupdate.go` SDK inside the binary

`selfupdate.go` SDK is very comprehensive, but the majority of the use cases can only call the single function at the beginning of the program.

```golang
package main

import (
    // ...
    "selfupdate.blockthrough.com"
    // ...
)


const (
    Version = "development"
)

func main() {
    selfupdate.Auto(
		context.Background(), // Context
		"blockthrough",       // Owner Name
		"selfupdate.go",      // Repo Name
		Version,              // Current Version
		"selfupdate",         // Executable Name
    )

    // rest of the program
}
```

`selfupdate.Auto` function automatically checks, downloads, patches and reruns the previously issued commands.

# Example

`selfupdate` cli is using itself for self-updating. Please refer to both `cmd/selfupdate/main.go` and `.github/workflows/build.yml` files for more info.
