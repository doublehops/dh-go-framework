package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIn(t *testing.T) {
	type fruit string

	fruits := []fruit{"apple", "banana", "carrot"}
	slice := make([]any, len(fruits))
	for i, f := range fruits {
		slice[i] = f
	}

	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		slice          []any
		value          any
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "stringFoundInSlice",
			errorMessage:   "",
			required:       true,
			slice:          slice,
			value:          "apple",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "stringNotFoundInSlice",
			errorMessage:   "",
			required:       true,
			slice:          slice,
			value:          "pineapple",
			expectedErrMsg: stringNotInSlice,
			expectedResult: false,
		},
		{
			name:           "stringEmptyButNotRequired",
			errorMessage:   "",
			required:       false,
			slice:          slice,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "stringEmptyButRequired",
			errorMessage:   "",
			required:       true,
			slice:          slice,
			value:          "",
			expectedErrMsg: RequiredPropertyError,
			expectedResult: false,
		},
		{
			name:           "customErrorMessage",
			errorMessage:   "not-found -- test",
			required:       true,
			slice:          slice,
			value:          "pineapple",
			expectedErrMsg: "not-found -- test",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := In(tt.slice, tt.errorMessage)
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}
