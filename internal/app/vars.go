// Package app provides shared application-level types and context keys.
package app

// ContextVar is a string type used as context keys to avoid collisions.
type ContextVar string

// UserIDKey and TraceIDKey are context keys for storing user ID and trace ID in request context.
const (
	UserIDKey  ContextVar = "userID"
	TraceIDKey ContextVar = "traceID"
)

func (c ContextVar) String() string {
	return string(c)
}
