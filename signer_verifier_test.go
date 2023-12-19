package selfupdate_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"selfupdate.blockthrough.com"
	"selfupdate.blockthrough.com/pkg/crypto"
)

func TestSignerVerifier(t *testing.T) {
	originalContent := "hello, world"
	originalContentReader := strings.NewReader(originalContent)

	publicKey, privateKey, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatal(err)
	}

	signer := selfupdate.NewHashSigner(privateKey)
	verifier := selfupdate.NewHashVerifier(publicKey)

	signedContentReader := signer.Sign(context.Background(), originalContentReader)
	verifiedContentReader := verifier.Verify(context.Background(), signedContentReader)

	verifiedContent, err := io.ReadAll(verifiedContentReader)
	if err != nil {
		t.Fatal(err)
	}

	if string(verifiedContent) != originalContent {
		t.Fatal("content is not matched")
	}
}
