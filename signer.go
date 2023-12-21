package selfupdate

import (
	"bytes"
	"context"
	"io"

	"selfupdate.blockthrough.com/pkg/crypto"
	"selfupdate.blockthrough.com/pkg/hash"
)

func NewHashSigner(privateKey crypto.PrivateKey) Signer {
	return SignerFunc(func(ctx context.Context, r io.Reader) io.Reader {
		var buffer bytes.Buffer

		hash, err := hash.FromReader(io.TeeReader(r, &buffer))
		if err != nil {
			return newErrorReader(err)
		}

		return io.MultiReader(
			bytes.NewReader(privateKey.Sign(hash)),
			&buffer,
		)
	})
}
