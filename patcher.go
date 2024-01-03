package selfupdate

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"

	"selfupdate.blockthrough.com/pkg/executil"
)

func NewPatcher(ext string) Patcher {
	return PatcherFunc(func(ctx context.Context, patch io.Reader) error {
		execPath, err := executil.CurrentPath()
		if err != nil {
			return err
		}

		target := execPath + ext

		err = writeToFile(target, patch)
		if err != nil {
			return err
		}

		if rc, ok := patch.(io.ReadCloser); ok {
			rc.Close()
		}

		// On Darwin, use the 'chmod' command to make the binary executable
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("chmod", "+x", target)
			return cmd.Run()
		}

		// For other platforms, use the 'os.Chmod' function
		return os.Chmod(target, 0755)
	})
}

func writeToFile(filename string, r io.Reader) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	return err
}
