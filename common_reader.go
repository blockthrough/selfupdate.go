package selfupdate

type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func newErrorReader(err error) *errorReader {
	return &errorReader{err}
}
