package selfupdate

import (
	"bytes"
	"context"
	"errors"
	"io"

	"selfupdate.blockthrough.com/pkg/crypto"
	"selfupdate.blockthrough.com/pkg/hash"
)

var (
	ErrVerificationFailed = errors.New("verification failed")
)

func NewHashVerifier(publicKey crypto.PublicKey) Verifier {
	return VerifierFunc(func(ctx context.Context, r io.Reader) io.Reader {
		var signedHash [hash.HashSize + crypto.Overhead]byte
		if _, err := io.ReadFull(r, signedHash[:]); err != nil {
			return newErrorReader(err)
		}

		var buffer bytes.Buffer

		contentHash, err := hash.FromReader(io.TeeReader(r, &buffer))
		if err != nil {
			return newErrorReader(err)
		}

		if !publicKey.Verify(signedHash[:]) {
			return newErrorReader(ErrVerificationFailed)
		}

		if !bytes.Equal(contentHash, signedHash[crypto.Overhead:]) {
			return newErrorReader(ErrVerificationFailed)
		}

		return &buffer
	})
}
