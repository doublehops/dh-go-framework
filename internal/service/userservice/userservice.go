// Package userservice provides business logic for user management.
package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/repository/userrepository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// UserService handles business logic for user management.
type UserService struct {
	*service.App
	userRepo *userrepository.Repo
}

// New creates a new UserService with the provided app and user repository.
func New(app *service.App, userRepo *userrepository.Repo) *UserService {
	return &UserService{
		App:      app,
		userRepo: userRepo,
	}
}

// Create persists a new user record after hashing the password.
func (s *UserService) Create(ctx context.Context, record *user.User) (*user.User, error) {
	record.OrganisationID = 1
	record.IsActive = 1
	hashedPassword, err := s.HashPassword(record.Password)
	if err != nil {
		s.Log.Error(ctx, "unable to has password", logga.KVPs{"error": err.Error()})

		return nil, errors.New("unable to has password")
	}

	record.Password = hashedPassword

	if err := record.SetCreated(ctx); err != nil {
		s.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err = s.userRepo.Create(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToSaveRecord+" "+err.Error(), nil)

		return record, req.ErrCouldNotSaveRecord
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	r := &user.User{}
	err = s.userRepo.GetByID(ctx, s.DB, record.ID, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return r, errors.New(service.UnableToRetrieveRecord)
	}

	return r, nil
}

// Update saves changes to an existing user record.
func (s *UserService) Update(ctx context.Context, record *user.User) (*user.User, error) {
	record.SetUpdated(ctx)

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.userRepo.Update(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToUpdateRecord+" "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	r := &user.User{}
	err = s.userRepo.GetByID(ctx, s.DB, record.ID, r)
	if err != nil {
		return r, errors.New(service.UnableToRetrieveRecord)
	}

	return r, nil
}

// DeleteByID soft-deletes a user record by setting its deleted_at timestamp.
func (s *UserService) DeleteByID(ctx context.Context, record *user.User) error {
	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	record.SetDeleted(ctx)

	err := s.userRepo.Delete(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, "unable to delete record. "+err.Error(), nil)

		return fmt.Errorf("unable to delete record")
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

// GetByID retrieves a user record by its ID.
func (s *UserService) GetByID(ctx context.Context, record *user.User, id int32) error {
	err := s.userRepo.GetByID(ctx, s.DB, id, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

// CreateUserSession creates a new authenticated session for the given user.
func (s *UserService) CreateUserSession(ctx context.Context, record *user.User) (*usersession.UserSession, error) {
	us := usersessionservice.New(s.App, nil)
	session, err := us.Create(ctx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)

		return nil, err
	}

	return session, nil
}

// GetAll retrieves a paginated collection of user records.
func (s *UserService) GetAll(ctx context.Context, r *req.Request) ([]*user.User, error) {
	records, err := s.userRepo.GetCollection(ctx, s.DB, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return records, nil
}

// GetByEmailAddress retrieves a user record by their email address.
func (s *UserService) GetByEmailAddress(ctx context.Context, record *user.User, emailAddress string) error {
	err := s.userRepo.GetByEmailAddress(ctx, s.DB, emailAddress, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}
