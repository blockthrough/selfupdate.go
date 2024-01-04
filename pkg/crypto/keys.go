package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/nacl/sign"
)

var (
	ErrInvalidKey = errors.New("invalid key")
)

const (
	Overhead       = sign.Overhead
	PublicKeySize  = 32
	PrivateKeySize = 64
)

type PrivateKey [PrivateKeySize]byte

func (p PrivateKey) Sign(message []byte) []byte {
	return sign.Sign(nil, message, (*[PrivateKeySize]byte)(&p))
}

func (p PrivateKey) String() string {
	return binary2String(p[:])
}

type PublicKey [PublicKeySize]byte

func (p PublicKey) Verify(signedMessage []byte) bool {
	_, ok := sign.Open(nil, signedMessage, (*[PublicKeySize]byte)(&p))
	return ok
}

func (p PublicKey) String() string {
	return binary2String(p[:])
}

func GenerateKeys() (publicKey PublicKey, privateKey PrivateKey, err error) {
	publicKeyPtr, privateKeyPtr, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	return *publicKeyPtr, *privateKeyPtr, nil
}

func ParsePrivateKey(key string) (priv PrivateKey, err error) {
	bytes, err := string2Binary(key)
	if err != nil {
		return
	}

	if len(bytes) != PrivateKeySize {
		return priv, ErrInvalidKey
	}

	copy(priv[:], bytes)
	return
}

func ParsePublicKey(key string) (pub PublicKey, err error) {
	bytes, err := string2Binary(key)
	if err != nil {
		return
	}

	if len(bytes) != PublicKeySize {
		return pub, ErrInvalidKey
	}

	copy(pub[:], bytes)
	return
}

// It seems like string2Binary is just a wrapper around hex.DecodeString.
// Same with binary2String and hex.EncodeToString.
// Do we actually need these functions?
// Or can we just use hex.DecodeString and hex.EncodeToString directly?
func string2Binary(str string) ([]byte, error) {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func binary2String(b []byte) string {
	return hex.EncodeToString(b)
}
