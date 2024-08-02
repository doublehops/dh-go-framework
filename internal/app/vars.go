package app

type ContextVar string

const (
	UserIDKey ContextVar = "userID"
)

func (c ContextVar) ToString() string {
	return string(c)
}
