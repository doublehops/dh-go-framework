package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinValue(t *testing.T) {
	tests := []struct {
		name           string
		minValue       int
		errorMessage   string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "goodValue",
			minValue:       20,
			errorMessage:   "",
			required:       true,
			value:          21,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "valueTooLow",
			minValue:       20,
			errorMessage:   "",
			required:       true,
			value:          10,
			expectedErrMsg: MinValueDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "valueTooLowCustomErrorMessage",
			minValue:       20,
			errorMessage:   "Value is too low -- test",
			required:       true,
			value:          10,
			expectedErrMsg: "Value is too low -- test",
			expectedResult: false,
		},
		{
			name:           "valueNotAnInt",
			minValue:       20,
			errorMessage:   "",
			required:       true,
			value:          "a",
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringAndRequired",
			minValue:       20,
			errorMessage:   "",
			required:       true,
			value:          "",
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringButNotRequired",
			minValue:       20,
			errorMessage:   "",
			required:       false,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := MinValue(tt.minValue, tt.errorMessage)
			res, validationErrMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, validationErrMsg, "error not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		name           string
		maxValue       int
		errorMessage   string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "goodValue",
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          10,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "valueTooHigh",
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          25,
			expectedErrMsg: MaxValueDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "valueTooHighCustomErrorMessage",
			maxValue:       20,
			errorMessage:   "Value is too low -- test",
			required:       true,
			value:          25,
			expectedErrMsg: "Value is too low -- test",
			expectedResult: false,
		},
		{
			name:           "valueNotAnInt",
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          "a",
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringAndRequired",
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          "",
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringButNotRequired",
			maxValue:       20,
			errorMessage:   "",
			required:       false,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := MaxValue(tt.maxValue, tt.errorMessage)
			res, validationErrMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, validationErrMsg, "error not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestIntInRange(t *testing.T) {
	tests := []struct {
		name           string
		minValue       int
		maxValue       int
		errorMessage   string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "goodValue",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          12,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "goodValueLowEnd",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          10,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "goodValueHighEnd",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          10,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "valueTooLow",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          2,
			expectedErrMsg: InRangeDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "valueTooHigh",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          25,
			expectedErrMsg: InRangeDefaultMessage,
			expectedResult: false,
		},
		{
			name:           "valueTooHighCustomErrorMessage",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "Value is too low -- test",
			required:       true,
			value:          25,
			expectedErrMsg: "Value is too low -- test",
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringAndRequired",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       true,
			value:          "",
			expectedErrMsg: ProcessingPropertyError,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringButNotRequired",
			minValue:       10,
			maxValue:       20,
			errorMessage:   "",
			required:       false,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := IntInRange(tt.minValue, tt.maxValue, tt.errorMessage)
			res, validationErrMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, validationErrMsg, "error not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestIsInt(t *testing.T) {
	tests := []struct {
		name           string
		errorMessage   string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{
			name:           "goodValue",
			errorMessage:   "",
			required:       true,
			value:          10,
			expectedErrMsg: "",
			expectedResult: true,
		},
		{
			name:           "badValue",
			errorMessage:   "",
			required:       true,
			value:          "a",
			expectedErrMsg: NotIntegerMessage,
			expectedResult: false,
		},
		{
			name:           "badValueWithCustomError",
			errorMessage:   "Not an int -- test",
			required:       true,
			value:          "a",
			expectedErrMsg: "Not an int -- test",
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringAndRequired",
			errorMessage:   "",
			required:       true,
			value:          "",
			expectedErrMsg: NotIntegerMessage,
			expectedResult: false,
		},
		{
			name:           "valueEmptyStringButNotRequired",
			errorMessage:   "",
			required:       false,
			value:          "",
			expectedErrMsg: "",
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := IsInt(tt.errorMessage)
			res, validationErrMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, validationErrMsg, "error not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}
