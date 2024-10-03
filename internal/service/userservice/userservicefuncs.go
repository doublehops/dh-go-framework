package userservice

import (
	"context"
	"errors"

	model "github.com/doublehops/dh-go-framework/internal/model/user"
)

// EmailAddressAlreadyExists will return bool for existence of email address.
func (s *UserService) EmailAddressAlreadyExists(ctx context.Context, emailAddress string) (bool, error) {
	record := &model.User{}
	if err := s.GetByEmailAddress(ctx, record, emailAddress); err != nil {
		s.Log.Error(ctx, "error retrieving record by email address"+err.Error(), nil)

		return false, errors.New("error retrieving record by email address")
	}

	if record.EmailAddress != "" {
		return true, nil
	}

	return false, nil
}
