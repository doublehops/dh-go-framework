// Package author defines the Author model and its validation rules.
package author

import (
	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

// Author is the database model for an author record.
type Author struct {
	model.BaseModel
	Name string `json:"name" db:"name"`
}

func (a *Author) getRules() []validator.Rule {
	return []validator.Rule{
		{"name", a.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 12, "")}}, //nolint:govet
	}
}

// Validate runs all validation rules for the author model.
func (a *Author) Validate() req.ErrMsgs {
	return validator.RunValidation(a.getRules())
}
