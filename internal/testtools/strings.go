package testtools

import "strings"

func Normalise(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), ""))
}
