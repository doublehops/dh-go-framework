// Package validator provides reusable field-level validation functions for models.
package validator

import "regexp"

// EmailAddressDefaultMessage is the default error message for invalid email addresses.
const (
	EmailAddressDefaultMessage = "is not a valid email address"

	emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// EmailAddress returns a ValidationFuncs that checks whether a value is a valid email address.
func EmailAddress(errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		var v string
		var ok bool

		if errorMessage == "" {
			errorMessage = EmailAddressDefaultMessage
		}

		if v, ok = value.(string); !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		re := regexp.MustCompile(emailPattern)
		if re.MatchString(v) {
			return true, ""
		}

		return false, errorMessage
	}
}
