package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailAddress(t *testing.T) {
	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "validEmailAddress",
			errorMessage:   "",
			required:       true,
			value:          "john@example.com",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "invalidEmailAddress",
			errorMessage:   "",
			required:       true,
			value:          "john-example.com",
			expectedErrMsg: EmailAddressDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "emailAddressEmptyButNotRequired",
			errorMessage:   "",
			required:       false,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "emailAddressEmptyButIsRequired",
			errorMessage:   "",
			required:       true,
			value:          "",
			expectedErrMsg: RequiredPropertyError,
			expectedResult: false,
		},
		{
			name:           "emailAddressNotAString",
			errorMessage:   "",
			required:       true,
			value:          nil,
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := EmailAddress(tt.errorMessage)
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}
