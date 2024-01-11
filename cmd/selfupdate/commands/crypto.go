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
		Usage: "work with public/private keys for signing and verifying",
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
		Usage: "genereating a pair of public/private keys for signing and verifying",
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
		Name:  "sign",
		Usage: "sign a binary using private key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Usage:    "content of the private key",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			key := ctx.String("key")

			privateKey, err := crypto.ParsePrivateKey(key)
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
		Name:  "verify",
		Usage: "verify a binary using public key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Usage:    "content of the public key",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			key := ctx.String("key")

			publicKey, err := crypto.ParsePublicKey(key)
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
