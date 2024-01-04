package selfupdate

import (
	"context"
	"os"
	"os/exec"
)

// NewCliRunner rerun the an executable with the same arguments.
// it requires the first argument to be the path to the executable.
func NewCliRunner(path string, args ...string) Runner {
	return RunnerFunc(func(ctx context.Context) error {
		cmd := exec.Command(path, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err := cmd.Start()
		if err != nil {
			return err
		}

		return cmd.Wait()
	})
}
