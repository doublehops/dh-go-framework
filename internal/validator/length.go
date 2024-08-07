package validator

const (
	MinLengthDefaultMessage     = "is below minimum required length"
	MaxLengthDefaultMessage     = "exceeds maximum length"
	BetweenLengthDefaultMessage = "is not within required range"
)

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
