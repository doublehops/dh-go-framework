// Package app provides shared application-level types and context keys.
package app

// ContextVar is a string type used as context keys to avoid collisions.
type ContextVar string

// UserIDKey, UserRecordKey and TraceIDKey are context keys for storing request-scoped values.
const (
	UserIDKey     ContextVar = "userID"
	UserRecordKey ContextVar = "userRecord"
	TraceIDKey    ContextVar = "traceID"
)

func (c ContextVar) String() string {
	return string(c)
}
