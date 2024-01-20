package util

import (
	"os"
	"path/filepath"
)

func ExecutionPath() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}
