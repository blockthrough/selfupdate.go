package compress

import (
	"compress/gzip"
	"io"
)

// Zip compresses data from the given io.Reader and returns a new io.Reader for the compressed data.
func Zip(r io.Reader) io.ReadCloser {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		var err error

		defer func() {
			if err != nil {
				pipeWriter.CloseWithError(err)
			} else {
				pipeWriter.Close()
			}
		}()

		writer := gzip.NewWriter(pipeWriter)
		defer writer.Close()

		_, err = io.Copy(writer, r)
	}()

	return pipeReader
}

// Unzip decompresses data from the given io.Reader and returns a new io.Reader for the decompressed data.
func Unzip(r io.Reader) io.ReadCloser {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		var err error

		defer func() {
			if err != nil {
				pipeWriter.CloseWithError(err)
			} else {
				pipeWriter.Close()
			}
		}()

		reader, err := gzip.NewReader(r)
		if err != nil {
			return
		}
		defer reader.Close()

		_, err = io.Copy(pipeWriter, reader)
	}()

	return pipeReader
}
