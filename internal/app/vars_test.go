package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		value    ContextVar
		expected string
	}{
		{
			name:     "expected",
			value:    UserIDKey,
			expected: string(UserIDKey),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.value.String())
		})
	}
}
