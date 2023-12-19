package hash

import (
	"crypto/sha256"
	"io"
)

const (
	HashSize = sha256.Size
)

func FromReader(r io.Reader) ([]byte, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
