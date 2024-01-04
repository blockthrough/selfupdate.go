package executil

import (
	"os"
	"path/filepath"
)

func CurrentPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// Clean up the path to get the absolute path without symbolic links
	return filepath.EvalSymlinks(exePath)
}
