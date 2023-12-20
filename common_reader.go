package selfupdate

import "io"

type errorReader struct {
	err error
}

var _ io.ReadCloser = (*errorReader)(nil)

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func (r *errorReader) Close() error {
	return nil
}

func newErrorReader(err error) *errorReader {
	return &errorReader{err}
}
