package selfupdate

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

		// On Darwin, use the 'chmod' command to make the binary executable
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("chmod", "+x", target)
			return cmd.Run()
		}

		// For other platforms, use the 'os.Chmod' function
		return os.Chmod(target, 0755)
	})
}
