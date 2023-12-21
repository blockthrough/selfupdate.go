package commands

import (
	"os"

	"selfupdate.blockthrough.com/pkg/cli"
	"selfupdate.blockthrough.com/pkg/crypto"
)

func generateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "geenerate values",
		Subcommands: []*cli.Command{
			generateKeys(),
		},
	}
}

func generateKeys() *cli.Command {
	return &cli.Command{
		Name:  "keys",
		Usage: "generate keys pair for signing and verifying",
		Action: func(ctx *cli.Context) error {
			publicKey, privateKey, err := crypto.GenerateKeys()
			if err != nil {
				return err
			}

			if err := createAndWrite("./selfupdate.pub", []byte(publicKey.String())); err != nil {
				return err
			}

			if err := createAndWrite("./selfupdate.key", []byte(privateKey.String())); err != nil {
				return err
			}

			return nil
		},
	}
}

func createAndWrite(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
