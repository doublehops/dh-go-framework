package author

import (
	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

type Author struct {
	model.BaseModel
	Name string `json:"name" db:"name"`
}

func (a *Author) getRules() []validator.Rule {
	return []validator.Rule{
		{"name", a.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 12, "")}}, //nolint:govet
	}
}

func (a *Author) Validate() req.ErrMsgs {
	return validator.RunValidation(a.getRules())
}
