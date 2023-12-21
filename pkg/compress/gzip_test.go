package compress_test

import (
	"io"
	"strings"
	"testing"

	"selfupdate.blockthrough.com/pkg/compress"
)

func TestZipUnzip(t *testing.T) {
	content := "hello, world"
	contentReader := strings.NewReader(content)

	zipReader := compress.Zip(contentReader)
	unzipReader := compress.Unzip(zipReader)

	unzippedContent, err := io.ReadAll(unzipReader)
	if err != nil {
		t.Fatal(err)
	}

	if string(unzippedContent) != content {
		t.Fatal("content is not matched")
	}
}
