package executil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CurrentPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Clean up the path to get the absolute path without symbolic links
	return filepath.EvalSymlinks(exePath)
}

func Copy(ext string) (newPath string, err error) {
	src, err := CurrentPath()
	if err != nil {
		return
	}

	if strings.HasSuffix(src, ext) {
		return "", nil
	}

	newPath = strings.TrimSuffix(src, ext)

	out, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	return
}

func Cleanup(ext string) error {
	src, err := CurrentPath()
	if err != nil {
		return err
	}

	if strings.HasSuffix(src, ext) {
		return nil
	}

	target := src + ext

	_, err = os.Stat(target)
	if err != nil {
		return nil
	}

	return os.Remove(target)
}
