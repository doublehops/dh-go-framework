package repository

import "github.com/google/uuid"

// NewID returns a new primary key: a UUIDv7 string for a CHAR(36) id column. v7 starts with a
// millisecond timestamp, so new ids sort after existing ones and InnoDB appends them to the end of
// the clustered primary-key index rather than scattering inserts across it as random v4s do. The
// timestamp is readable from the id, so it reveals when the record was created.
//
// NewV7 only fails if crypto/rand does, which on supported platforms aborts the process rather
// than returning an error, so it is treated as unrecoverable — as uuid.New() already did.
func NewID() string {
	return uuid.Must(uuid.NewV7()).String()
}
