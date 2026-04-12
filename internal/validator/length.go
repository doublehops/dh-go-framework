package validator

// MinLengthDefaultMessage and related constants are default error messages for string length validation.
const (
	MinLengthDefaultMessage     = "is below minimum required length"
	MaxLengthDefaultMessage     = "exceeds maximum length"
	BetweenLengthDefaultMessage = "length is not within required range"
)

// MinLength returns a ValidationFuncs that checks the string meets a minimum length requirement.
func MinLength(minLength int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = MinLengthDefaultMessage
		}

		var v string
		var ok bool

		if v, ok = value.(string); !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		if len(v) < minLength {
			return false, errorMessage
		}

		return true, ""
	}
}

// MaxLength returns a ValidationFuncs that checks the string does not exceed a maximum length.
func MaxLength(maxLength int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = MaxLengthDefaultMessage
		}

		var v string
		var ok bool

		if v, ok = value.(string); !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		if len(v) > maxLength {
			return false, errorMessage
		}

		return true, ""
	}
}

// LengthInRange returns a ValidationFuncs that checks the string length is within a min/max range.
func LengthInRange(minLength, maxLength int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = BetweenLengthDefaultMessage
		}

		var v string
		var ok bool

		if v, ok = value.(string); !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		if len(v) >= minLength && len(v) <= maxLength {
			return true, ""
		}

		return false, errorMessage
	}
}
