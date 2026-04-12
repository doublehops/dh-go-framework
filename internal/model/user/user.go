// Package user defines the User model, request/response types, and validation rules.
package user

import (
	"time"

	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

// User is the full database model for a user record.
type User struct {
	model.BaseModel
	OrganisationID      int        `json:"organisationId" db:"organisation_id"`
	Name                string     `json:"name" db:"name"`
	EmailAddress        string     `json:"emailAddress" db:"email_address"`
	EmailVerified       bool       `json:"emailVerified" db:"email_verified"`
	EmailVerifiedToken  string     `json:"emailVerifiedToken" db:"email_verified_token"`
	Password            string     `json:"password" db:"password"`
	PasswordResetCode   string     `json:"passwordResetToken" db:"password_reset_token"`
	PasswordResetExpire *time.Time `json:"passwordResetExpire" db:"password_reset_expire"`
	IsActive            int        `json:"isActive" db:"is_active"`
}

// ResponseUser is a pruned user model for API responses, omitting sensitive fields.
type ResponseUser struct {
	model.BaseModel
	OrganisationID int    `json:"organisationId"`
	Name           string `json:"name"`
	EmailAddress   string `json:"emailAddress"`
	EmailVerified  bool   `json:"emailVerified"`
	IsActive       int    `json:"isActive"`
}

// CreateUser holds the fields required for creating a new user.
type CreateUser struct {
	model.BaseModel
	Name         string `json:"name" db:"name"`
	EmailAddress string `json:"emailAddress" db:"email_address"`
	Password     string `json:"password" db:"password"`
}

// SuccessfulLoginResponse is the response payload returned after a successful login.
type SuccessfulLoginResponse struct {
	BearerToken string `json:"bearerToken"`
}

// UpdateUser holds the fields allowed when updating a user record.
type UpdateUser struct {
	model.BaseModel
	Name string `json:"name" db:"name"`
}

func (u *User) getUserCreateRules() []validator.Rule {
	return []validator.Rule{
		{"name", u.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},          //nolint:govet
		{"emailAddress", u.EmailAddress, true, []validator.ValidationFuncs{validator.EmailAddress("")}}, //nolint:govet
		{"password", u.Password, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},  //nolint:govet

	}
}

func (u *User) getUserUpdateRules() []validator.Rule {
	return []validator.Rule{
		{"name", u.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}}, //nolint:govet
	}
}

// ValidateCreate runs validation rules for user creation requests.
func (u *User) ValidateCreate() req.ErrMsgs {
	return validator.RunValidation(u.getUserCreateRules())
}

// ValidateUpdate runs validation rules for user update requests.
func (u *User) ValidateUpdate() req.ErrMsgs {
	return validator.RunValidation(u.getUserUpdateRules())
}

// ValidateLogin runs validation rules for login requests.
func (u *User) ValidateLogin() req.ErrMsgs {
	return validator.RunValidation([]validator.Rule{
		// {"organisationId", u.OrganisationID, true, []validator.ValidationFuncs{validator.IsInt("")}},                         //nolint:govet
		{"emailAddress", u.EmailAddress, true, []validator.ValidationFuncs{validator.EmailAddress("")}}, //nolint:govet
		{"password", u.Password, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},  //nolint:govet
	})
}
