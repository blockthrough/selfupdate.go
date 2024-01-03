package commands

import (
	"context"
	"io"
	"os"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/pkg/cli"
	"selfupdate.blockthrough.com/pkg/crypto"
)

func cryptoCmd() *cli.Command {
	return &cli.Command{
		Name:  "crypto",
		Usage: "geenerate values",
		Subcommands: []*cli.Command{
			cryptoGenerateKeys(),
			cryptoSign(),
			cryptoVerify(),
		},
	}
}

func cryptoGenerateKeys() *cli.Command {
	return &cli.Command{
		Name:  "keys",
		Usage: "genereating crypto keys pair for signing and verifying",
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

func cryptoSign() *cli.Command {
	return &cli.Command{
		Name:        "sign",
		Usage:       "sign a binary using private key",
		Description: `make sure to set SELF_UPDATE_PRIVATE_KEY env variable before using this command`,
		Action: func(ctx *cli.Context) error {
			privateKey, err := getPrivateKey()
			if err != nil {
				return err
			}

			signer := selfupdate.NewHashSigner(privateKey)
			_, err = io.Copy(os.Stdout, signer.Sign(context.Background(), os.Stdin))
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func cryptoVerify() *cli.Command {
	return &cli.Command{
		Name:        "verify",
		Usage:       "verify a binary using public key",
		Description: `make sure to set SELF_UPDATE_PUBLIC_KEY env variable before using this command`,
		Action: func(ctx *cli.Context) error {
			publicKey, err := getPublicKey()
			if err != nil {
				return err
			}

			verifier := selfupdate.NewHashVerifier(publicKey)
			_, err = io.Copy(os.Stdout, verifier.Verify(context.Background(), os.Stdin))
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func getPublicKey() (publicKey crypto.PublicKey, err error) {
	value, ok := os.LookupEnv("SELF_UPDATE_PUBLIC_KEY")
	if !ok || value == "" {
		err = cli.Exit("SELF_UPDATE_PUBLIC_KEY env variable is not set", 1)
		return
	}

	publicKey, err = crypto.ParsePublicKey(value)
	return
}

func getPrivateKey() (privateKey crypto.PrivateKey, err error) {
	value, ok := os.LookupEnv("SELF_UPDATE_PRIVATE_KEY")
	if !ok || value == "" {
		err = cli.Exit("SELF_UPDATE_PRIVATE_KEY env variable is not set", 1)
		return
	}

	privateKey, err = crypto.ParsePrivateKey(value)
	return
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
