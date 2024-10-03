package user

import (
	"time"

	"github.com/doublehops/dh-go-framework/internal/model"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/validator"
)

type User struct {
	model.BaseModel
	OrganisationID      int        `json:"organisationId" db:"organisation_id"`
	Name                string     `json:"name" db:"name"`
	EmailAddress        string     `json:"emailAddress" db:"email_address"`
	EmailVerified       bool       `json:"emailVerified" db:"email_verified"`
	EmailVerifiedToken  int        `json:"emailVerifiedToken" db:"email_verified_token"`
	Password            string     `json:"password" db:"password"`
	PasswordResetCode   string     `json:"passwordResetToken" db:"password_reset_token"`
	PasswordResetExpire *time.Time `json:"passwordResetExpire" db:"password_reset_expire"`
	IsActive            int        `json:"isActive" db:"is_active"`
}

func (u *User) getRules() []validator.Rule {
	return []validator.Rule{
		// {"organisationId", u.OrganisationID, true, []validator.ValidationFuncs{validator.IsInt("")}},                         //nolint:govet
		{"name", u.Name, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},          //nolint:govet
		{"emailAddress", u.EmailAddress, true, []validator.ValidationFuncs{validator.EmailAddress("")}}, //nolint:govet
		// {"emailVerified", u.EmailVerified, true, []validator.ValidationFuncs{validator.IsInt("")}},                           //nolint:govet
		// {"password", u.Password, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}},                       //nolint:govet
		// {"passwordResetString", u.PasswordResetString, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}}, //nolint:govet
		// {"passwordResetExpire", u.PasswordResetExpire, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, "")}}, //nolint:govet,
		// {"isActive", u.IsActive, true, []validator.ValidationFuncs{validator.IsInt("")}},                                     //nolint:govet

	}
}

func (u *User) Validate() req.ErrMsgs {
	return validator.RunValidation(u.getRules())
}
