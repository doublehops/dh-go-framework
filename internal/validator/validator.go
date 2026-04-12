package validator

import (
	req "github.com/doublehops/dh-go-framework/internal/request"
)

// RequiredPropertyError and ProcessingPropertyError are common validation error messages.
const (
	RequiredPropertyError   = "this is a required property"
	ProcessingPropertyError = "unable to process property"
)

// ValidationFuncs is a function type that takes a required flag and value, and returns a pass/fail and message.
type ValidationFuncs func(bool, interface{}) (bool, string)

// Rule defines a single validation rule with a field name, value, required flag, and functions to run.
type Rule struct {
	VariableName string
	Value        interface{}
	Required     bool
	Function     []ValidationFuncs
}

// Error is a string type for validation error messages.
type Error string

// RunValidation executes all provided validation rules and returns a map of field errors.
func RunValidation(rules []Rule) req.ErrMsgs {
	errorMessages := make(req.ErrMsgs)

	for _, prop := range rules {
		var errors []string
		for _, rule := range prop.Function {
			valid, errMsg := rule(prop.Required, prop.Value)
			if !valid {
				errors = append(errors, errMsg)
			}
		}

		if errors != nil {
			errorMessages[prop.VariableName] = errors
		}
	}

	return errorMessages
}
