// Package tools provides general utility functions used across the application.
package tools

import (
	"runtime"
)

// CurrentFunction will return the function two calls back to be added to log messages.
func CurrentFunction() string {
	counter, _, _, success := runtime.Caller(2)

	if !success {
		return ""
	}

	c := runtime.FuncForPC(counter)

	if c == nil {
		return "<unknown>"
	}

	return c.Name()
}
