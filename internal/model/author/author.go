package author

import (
	"github.com/doublehops/dh-go-framework/internal/model"
	"github.com/doublehops/dh-go-framework/internal/validator"

	req "github.com/doublehops/dh-go-framework/internal/request"
)

type Author struct {
	model.BaseModel
	Name string `json:"name"`
}

func (a *Author) getRules() []validator.Rule {
	return []validator.Rule{
		{"name", a.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}}, //nolint:govet
	}
}

func (a *Author) Validate() req.ErrMsgs {
	return validator.RunValidation(a.getRules())
}
