package crypto_test

import (
	"bytes"
	"testing"

	"selfupdate.blockthrough.com/pkg/crypto"
)

func TestSignVerify(t *testing.T) {
	publicKey, privateKey, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("hello, world")
	signedMessage := privateKey.Sign(message)

	if len(signedMessage) != len(message)+crypto.Overhead {
		t.Fatal("sign size overhead is not matched")
	}

	ok := publicKey.Verify(signedMessage)
	if !ok {
		t.Fatal("verify failed")
	}
}

func TestEncodeDecodePublicKey(t *testing.T) {
	publicKey, _, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	strValue := publicKey.String()
	key, err := crypto.ParsePublicKey(strValue)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(publicKey[:], key[:]) {
		t.Fatal("public key is not matched")
	}
}

func TestEncodeDecodePrivateKey(t *testing.T) {
	_, privateKey, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	strValue := privateKey.String()
	key, err := crypto.ParsePrivateKey(strValue)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(privateKey[:], key[:]) {
		t.Fatal("public key is not matched")
	}
}
