package app

type ContextVar string

const (
	UserIDKey  ContextVar = "userID"
	TraceIDKey ContextVar = "traceID"
)

func (c ContextVar) String() string {
	return string(c)
}
