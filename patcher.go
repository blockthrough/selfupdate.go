package selfupdate

import (
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
)

func NewPatcher(outfile string) Patcher {
	return PatcherFunc(func(ctx context.Context, patch io.Reader) error {
		err := writeToFile(outfile, patch)
		if err != nil {
			return err
		}

		// On Darwin, use the 'chmod' command to make the binary executable
		if runtime.GOOS == "darwin" {
			cmd := exec.Command("chmod", "+x", outfile)
			return cmd.Run()
		}

		// For other platforms, use the 'os.Chmod' function
		return os.Chmod(outfile, 0755)
	})
}

func writeToFile(filename string, r io.Reader) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}

	_, err = io.Copy(out, r)

	return err
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	return writeToFile(dst, in)
}
