package selfupdate

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

func NewPatcher() Patcher {
	return PatcherFunc(func(ctx context.Context, patch io.Reader) error {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}

		// Clean up the path to get the absolute path without symbolic links
		target, err := filepath.EvalSymlinks(exePath)
		if err != nil {
			return err
		}

		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, patch)
		if err != nil {
			return err
		}

		return nil
	})
}
