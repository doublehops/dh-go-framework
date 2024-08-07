package validator

import (
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunValidation(t *testing.T) {
	type fruit string

	tests := []struct {
		name               string
		rules              []Rule
		expectedErrorCount int
		expectedErrMsgs    req.ErrMsgs
		expectedResult     bool
	}{
		{
			name: "validationSuccess",
			rules: []Rule{
				{"name", "john", true, []ValidationFuncs{LengthInRange(3, 4, "")}}, //nolint:govet
			},
			expectedErrorCount: 0,
			expectedErrMsgs:    nil,
			expectedResult:     true,
		},
		{
			name: "validationFail",
			rules: []Rule{
				{"name", "johnathon", true, []ValidationFuncs{LengthInRange(3, 4, "")}}, //nolint:govet
			},
			expectedErrorCount: 1,
			expectedErrMsgs: req.ErrMsgs{
				"name": []string{InRangeDefaultMessage},
			},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := RunValidation(tt.rules)
			assert.Equal(t, tt.expectedErrorCount, len(res))
		})
	}
}
