package userservice

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

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

// HashPassword hashes a password using bcrypt.
func (s *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash checks if the hashed password matches the plaintext password.
func (s *UserService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
