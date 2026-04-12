// Package helpers provides utility functions for the migration package.
package helpers

import "os"

// PrintMsg writes a message to stderr.
func PrintMsg(msg string) {
	_, _ = os.Stderr.WriteString(msg)
}
