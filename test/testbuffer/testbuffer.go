// Package testbuffer provides a file-backed io.Writer used to capture log output in tests.
package testbuffer

import (
	"fmt"
	"os"
)

// Filename is the path to the temporary log file used by the test buffer.
var Filename = "/tmp/testwriter.log"

// TestBuffer implements io.Writer by writing log output to a temp file for test assertions.
type TestBuffer struct{}

func (tb TestBuffer) Read() ([]byte, error) {
	data, err := os.ReadFile(Filename)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read file: %s. %s", Filename, err)
	}

	err = os.Remove(Filename)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to remove file: %s. %s", Filename, err)
	}

	return data, nil
}

// Write will write contents to a file. This will not work is tests are running asynchronous.
func (tb TestBuffer) Write(p []byte) (n int, err error) {
	err = os.WriteFile(Filename, p, 0o600)
	if err != nil {
		return 0, fmt.Errorf("unable to write temp file: %s. %s", Filename, err)
	}

	return len(p), nil
}
