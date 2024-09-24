package user

import (
	"time"

	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

type User struct {
	model.BaseModel
	OrganisationID       int        `json:"organisationId"`
	Name                 string     `json:"name"`
	EmailAddress         string     `json:"emailAddress"`
	EmailVerified        int        `json:"emailVerified"`
	Password             string     `json:"password"`
	PasswordResetString  string     `json:"passwordResetString"`
	PasswordResetTimeout *time.Time `json:"passwordResetTimeout"`
	IsActive             int        `json:"isActive"`
}

func (u *User) getRules() []validator.Rule {
	return []validator.Rule{
		{"organisationId", u.OrganisationID, true, []validator.ValidationFuncs{validator.IsInt("")}},                           //nolint:govet
		{"name", u.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},                                 //nolint:govet
		{"emailAddress", u.EmailAddress, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},                 //nolint:govet
		{"emailVerified", u.EmailVerified, true, []validator.ValidationFuncs{validator.IsInt("")}},                             //nolint:govet
		{"password", u.Password, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},                         //nolint:govet
		{"passwordResetString", u.PasswordResetString, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},   //nolint:govet
		{"passwordResetTimeout", u.PasswordResetTimeout, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}}, //nolint:govet,
		{"isActive", u.IsActive, true, []validator.ValidationFuncs{validator.IsInt("")}},                                       //nolint:govet

	}
}

func (u *User) Validate() req.ErrMsgs {
	return validator.RunValidation(u.getRules())
}
