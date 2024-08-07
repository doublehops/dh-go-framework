package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinLength(t *testing.T) {
	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		minLength      int
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "validLength",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			value:          "apple",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "invalidLength",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			value:          "frog",
			expectedErrMsg: MinLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "customErrorMessage",
			errorMessage:   "Not valid -- test",
			required:       true,
			minLength:      5,
			value:          "frog",
			expectedErrMsg: "Not valid -- test",
			expectedResult: false,
		},
		{
			name:           "invalidNotString",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			value:          nil,
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "invalidLengthAndNotRequired",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			value:          "frog",
			expectedErrMsg: MinLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "emptyStringNotRequired",
			errorMessage:   "",
			required:       false,
			minLength:      5,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "emptyStringButRequired",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			value:          "",
			expectedErrMsg: RequiredPropertyError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := MinLength(tt.minLength, tt.errorMessage)
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		maxLength      int
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "validLength",
			errorMessage:   "",
			required:       true,
			maxLength:      5,
			value:          "apple",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "invalidLength",
			errorMessage:   "",
			required:       true,
			maxLength:      5,
			value:          "banana",
			expectedErrMsg: MaxLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "customErrorMessage",
			errorMessage:   "Not valid -- test",
			required:       true,
			maxLength:      5,
			value:          "banana",
			expectedErrMsg: "Not valid -- test",
			expectedResult: false,
		},
		{
			name:           "invalidNotString",
			errorMessage:   "",
			required:       true,
			maxLength:      5,
			value:          nil,
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "invalidLengthAndNotRequired",
			errorMessage:   "",
			required:       true,
			maxLength:      5,
			value:          "banana",
			expectedErrMsg: MaxLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "emptyStringNotRequired",
			errorMessage:   "",
			required:       false,
			maxLength:      5,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "emptyStringButRequired",
			errorMessage:   "",
			required:       true,
			maxLength:      5,
			value:          "",
			expectedErrMsg: RequiredPropertyError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := MaxLength(tt.maxLength, tt.errorMessage)
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestLengthInRange(t *testing.T) {
	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		minLength      int
		maxLength      int
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "validLength",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          "apple",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "invalidLengthTooShort",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          "pie",
			expectedErrMsg: BetweenLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "invalidLengthTooLong",
			errorMessage:   "",
			required:       true,
			minLength:      1,
			maxLength:      5,
			value:          "banana",
			expectedErrMsg: BetweenLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "customErrorMessage",
			errorMessage:   "Not valid -- test",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          "pie",
			expectedErrMsg: "Not valid -- test",
			expectedResult: false,
		},
		{
			name:           "invalidNotString",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          nil,
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "invalidLengthAndNotRequired",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          "pie",
			expectedErrMsg: BetweenLengthDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "emptyStringNotRequired",
			errorMessage:   "",
			required:       false,
			minLength:      5,
			maxLength:      15,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "emptyStringButRequired",
			errorMessage:   "",
			required:       true,
			minLength:      5,
			maxLength:      15,
			value:          "",
			expectedErrMsg: RequiredPropertyError,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := LengthInRange(tt.minLength, tt.maxLength, tt.errorMessage)
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}
