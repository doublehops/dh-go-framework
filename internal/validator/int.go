package validator

// MinValueDefaultMessage and related constants are default error messages for integer validation functions.
const (
	MinValueDefaultMessage = "is below required amount"
	MaxValueDefaultMessage = "is above required amount"
	InRangeDefaultMessage  = "is not within required range"
	NotIntegerMessage      = "is not an integer"
)

// MinValue returns a ValidationFuncs that checks the value is at or above a minimum integer.
func MinValue(minValue int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = MinValueDefaultMessage
		}

		var (
			v  int
			ok bool
		)

		if value == "" && !required {
			return true, ""
		}

		if v, ok = value.(int); !ok {
			return false, ProcessingPropertyError
		}

		if v < minValue {
			return false, errorMessage
		}

		return true, ""
	}
}

// MaxValue returns a ValidationFuncs that checks the value does not exceed a maximum integer.
func MaxValue(maxValue int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = MaxValueDefaultMessage
		}

		var (
			v  int
			ok bool
		)

		if value == "" && !required {
			return true, ""
		}

		if v, ok = value.(int); !ok {
			return false, ProcessingPropertyError
		}

		if v > maxValue {
			return false, errorMessage
		}

		return true, ""
	}
}

// IntInRange returns a ValidationFuncs that checks the value is within a min/max integer range.
func IntInRange(minValue, maxValue int, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = InRangeDefaultMessage
		}

		var (
			v  int
			ok bool
		)

		if value == "" && !required {
			return true, ""
		}

		if v, ok = value.(int); !ok {
			return false, ProcessingPropertyError
		}

		if v < minValue || v > maxValue {
			return false, errorMessage
		}

		return true, ""
	}
}

// IsInt returns a ValidationFuncs that checks whether the value is an integer.
func IsInt(errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = NotIntegerMessage
		}

		if value == "" && !required {
			return true, ""
		}

		var ok bool

		if _, ok = value.(int); !ok {
			return false, errorMessage
		}

		return true, ""
	}
}
