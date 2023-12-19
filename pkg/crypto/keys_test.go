package crypto

import "testing"

func TestSignVerify(t *testing.T) {
	publicKey, privateKey, err := GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("hello, world")
	signedMessage := privateKey.Sign(message)

	if len(signedMessage) != len(message)+Overhead {
		t.Fatal("sign size overhead is not matched")
	}

	ok := publicKey.Verify(signedMessage)
	if !ok {
		t.Fatal("verify failed")
	}
}
