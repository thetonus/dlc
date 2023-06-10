package fileutils

import (
	"io"
	"os"
)

// ReadFile reads a YAML file
func ReadFile(path string) ([]byte, error) {
	// Handle /dev/stdin
	if path == "-" || path == "/dev/stdin" {
		return io.ReadAll(os.Stdin)
	}
	// Handle normal files
	return os.ReadFile(path)
}
