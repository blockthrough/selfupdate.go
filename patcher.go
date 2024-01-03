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

		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, patch)
		if err != nil {
			return err
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
