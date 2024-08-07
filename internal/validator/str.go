package validator

import "fmt"

const (
	stringNotInSlice = "value was not found in available items"
)

// nolint:cyclop
func In(slice []any, errorMessage string) ValidationFuncs {
	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = stringNotInSlice
		}

		// We want to convert the slice values into strings.
		var strSlice []string
		for _, i := range slice {
			item := fmt.Sprintf("%s", i)
			strSlice = append(strSlice, item)
		}

		v, ok := value.(string)
		if !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		for _, item := range strSlice {
			if item == v {
				return true, ""
			}
		}

		return false, errorMessage
	}
}
