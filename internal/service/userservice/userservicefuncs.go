package userservice

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// EmailAddressAlreadyExists will return bool for existence of email address.
func (s *UserService) EmailAddressAlreadyExists(ctx context.Context, emailAddress string) (bool, error) {
	record, err := s.GetByEmailAddress(ctx, emailAddress)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		s.Log.Error(ctx, "error retrieving record by email address"+err.Error(), nil)

		return false, errors.New("error retrieving record by email address")
	}

	return record.Record.EmailAddress != "", nil
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
