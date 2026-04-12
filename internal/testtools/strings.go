package testtools

import "strings"

// Normalise lowercases and removes whitespace from a string for comparison in tests.
func Normalise(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), ""))
}
