package selfupdate

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
)

// NewCliRunner rerun the an executable with the same arguments.
// it requires the first argument to be the path to the executable.
func NewCliRunner(path string, args ...string) Runner {
	return RunnerFunc(func(ctx context.Context) error {
		cmd := exec.Command(path, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Start()
	})
}

func NewAutoCliRunner() Runner {
	return RunnerFunc(func(ctx context.Context) error {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}

		// Clean up the path to get the absolute path without symbolic links
		target, err := filepath.EvalSymlinks(exePath)
		if err != nil {
			return err
		}

		return NewCliRunner(target, os.Args[1:]...).Run(ctx)
	})
}
