package fileutils

import (
	"bufio"
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

func WriteFile(f io.Writer, content []byte) error {
	writer := bufio.NewWriter(f)
	_, err := writer.Write(content)
	defer writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
