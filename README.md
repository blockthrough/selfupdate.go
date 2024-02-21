# Self Update for Golang

It's a dead-simple toolchain for updating Golang binaries. It is designed to leverage Github's Releases and Actions. It comes with its CLI and SDK to make things even simpler.

## CLI Installation

```bash
go install selfupdate.blockthough.com/cmd/selfupdate@latest
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

sign a binary using a private key. The private key must be passed as an argument using `--key`

```bash
selfupdate crypto sign --key "CONTENT OF PRIVATE KEY" < ./bin/file > ./bin/file.sig
```

#### verify

verify a binary using a public key. The public key must be passed as an argument using `--key`

```bash
selfupdate crypto verify --key "CONTENT OF PUBLIC KEY"  < ./bin/file.sg > ./bin/file
```

### github

a provider tool for working with Github's apis for releasing, uploading and downloading binaries.

> NOTE: an environment variable, `SELF_UPDATE_GH_TOKEN`, needs to be created and set with `GITHUB_TOKEN`. It is highly recommened to use it `GITHUB_TOKEN` rather than personal token. Also, this command with all its sub commands needs to be called inside github's actions workflow.

#### check

check if there is a new version by providing the version

```bash
selfupdate github check --owner blockthough --repo selfupdate.go --filename selfupdate --version v0.0.1
```

> for more info, run `selfupdate github check --help`

#### release

To provide a better developer experience, this command can create a new Github release, with a title, and a description and attach it to a specific tag.

> NOTE: running this command with the same version will be noop. So it can be safely used in github actions workflow with matrix strategy.

```bash
selfupdate github release -owner blockthough --repo selfupdate.go --token GITHUB_TOKEN --version v0.0.1 --title version v0.0.1 --desc "this is an amazin release"
```

> for more info, run `selfupdate github release --help`

#### upload

upload a new asset to an already created Github release. If required to sign the binary file before upload, provide `--key` with the content of the generated private key.

> In order to upload assets, a github release must be created first. Please refer to `release` subcommand. Also this command can be used multiple times for each individual asset in github actions workflow.

```bash
selfupdate github upload -owner blockthough --repo selfupdate.go --token GITHUB_TOKEN --version v0.0.1 --filename selfupload.sign --key PRIVATE_KEY < /path/to/file
```

#### download

To download a specific asset from Github releases, this command can be used. It requires `--filename` and `--version` to be presented.

```bash
selfupdate github download -owner blockthough --repo selfupdate.go --version v0.0.1 --filename selfupload.sign --key PUBLIC_KEY > /path/to/file
```

## Usage

To have successful self-updating binaries, two steps need to be followed:

### ( 1 ) Github Actions Workflow

First, compile your code and generate a binary, make sure to use `-ldflags "-X main.Version=${{ github.ref_name }} -X main.PublicKey=${{ PUBLIC_KEY }}"` flag during `go build` to inject the new tag as a version into the binary.

- Create a new Release using `selfupdate github release` command
- Sign and upload content using `selfupdate github upload`

### ( 2 ) Using `selfudpate` SDK inside the binary

inside the go project install the SKD by running the following command:

```bash
go get selfupdate.blockthrough.com@latest
```

`selfupdate.go` SDK is very comprehensive, but the majority of the projects can only call the single function at the beginning of the program.

```golang
package main

import (
    // ...
    "selfupdate.blockthrough.com"
    // ...
)


const (
    Version = ""
    PublicKey = ""
)

func main() {
    // NOTE: please refer to "Create a Fine-Grained Personal Access Tokens" section of doc
    ghToken, ok := os.LookupEnv("MY_AWESOME_PROJECT_GITHUB_TOKEN")
    if !ok {
        // error out that github token is not presented
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

    // rest of the program
}
```

`selfupdate.Auto` function automatically checks, downloads, patches and re-runs the previously issued command.

# Example

`selfupdate` CLI is using itself for self-updating. Please refer to both `cmd/selfupdate/main.go` and `.github/workflows/build.yml` files for more info.

# Create Fine-Grained Personal Access Tokens

Each person who needs to use your app CLI and leverage the self-updating is required to create a GitHub API token. It is recommended to use the following [Token/Settings](https://github.com/settings/tokens?type=beta) to generate the API keys.

The only required option is to select the project, and on Repository Permissions, select only Contents as Read access. We only need to read the metadata and download assets during the update, nothing more.
