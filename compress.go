package selfupdate

import (
	"compress/gzip"
	"io"
)

func compressReader(r io.Reader) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		var err error
		var zw *gzip.Writer

		defer func() {
			if err != nil {
				zw.Close()
				pw.CloseWithError(err)
			} else {
				zw.Close()
				pw.Close()
			}

			if rc := r.(io.ReadCloser); rc != nil {
				rc.Close()
			}
		}()

		zw = gzip.NewWriter(pw)
		_, err = io.Copy(zw, r)
	}()

	return pr
}

func decompressReader(r io.Reader) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		var err error
		var zr *gzip.Reader

		defer func() {
			if err != nil {
				zr.Close()
				pw.CloseWithError(err)
			} else {
				zr.Close()
				pw.Close()
			}

			if rc := r.(io.ReadCloser); rc != nil {
				rc.Close()
			}
		}()

		zr, err = gzip.NewReader(pr)
		if err != nil {
			return
		}

		_, err = io.Copy(pw, zr)
	}()

	return pr
}
