// Package author defines the Author model and its validation rules.
package author

import (
	"strings"

	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

// Author is the database model for an author record.
type Author struct {
	model.BaseModel
	Name string `json:"name" db:"name"`
}

// Normalise tidies user input before validation, so the length rules apply to the stored value
// (an all-whitespace name correctly fails the minimum-length check). Call before Validate.
func (a *Author) Normalise() {
	a.Name = strings.TrimSpace(a.Name)
}

func (a *Author) getRules() []validator.Rule {
	return []validator.Rule{
		{"name", a.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 32, "")}}, //nolint:govet
	}
}

// Validate runs all validation rules for the author model.
func (a *Author) Validate() req.ErrMsgs {
	return validator.RunValidation(a.getRules())
}
