package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/sign"
)

const (
	Overhead = sign.Overhead
)

type PrivateKey [64]byte

func (p PrivateKey) Sign(message []byte) []byte {
	return sign.Sign(nil, message, (*[64]byte)(&p))
}

type PublicKey [32]byte

func (p PublicKey) Verify(signedMessage []byte) bool {
	_, ok := sign.Open(nil, signedMessage, (*[32]byte)(&p))
	return ok
}

func GenerateKeys() (publicKey PublicKey, privateKey PrivateKey, err error) {
	publicKeyPtr, privateKeyPtr, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	return *publicKeyPtr, *privateKeyPtr, nil
}
