package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	first := NewID()
	second := NewID()

	parsed, err := uuid.Parse(first)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, uuid.Version(7), parsed.Version())
	assert.Len(t, first, 36, "must fit a CHAR(36) column")

	// v7 ids are time-ordered, so a later id sorts after an earlier one as a string too.
	assert.Less(t, first, second)
}
