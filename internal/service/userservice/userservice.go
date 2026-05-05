// Package userservice provides business logic for user management.
package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/repository/userrepository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"
)

// UserRepository defines the data access methods required by UserService.
type UserRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	Update(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	Delete(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	GetByID(ctx context.Context, db *sqlx.DB, id int32, record *user.User) error
	GetByEmailAddress(ctx context.Context, db *sqlx.DB, emailAddress string, record *user.User) error
	GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*user.User, error)
}

// UserService handles business logic for user management.
type UserService struct {
	*service.App
	userRepo UserRepository
}

// New creates a new UserService, wiring its own repository internally.
func New(app *service.App) *UserService {
	return &UserService{
		App:      app,
		userRepo: userrepository.New(app.Log),
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
	us := usersessionservice.New(s.App)
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

// DeleteSessionByToken looks up a session by token and soft-deletes it.
func (s *UserService) DeleteSessionByToken(ctx context.Context, token string) error {
	us := usersessionservice.New(s.App)

	record := &usersession.UserSession{}
	if err := us.GetByToken(ctx, record, token); err != nil || record.ID == 0 {
		return errors.New(service.UnableToRetrieveRecord)
	}

	return us.DeleteByID(ctx, record)
}

// GetByEmailAddress retrieves a user record by their email address.
func (s *UserService) GetByEmailAddress(ctx context.Context, record *user.User, emailAddress string) error {
	err := s.userRepo.GetByEmailAddress(ctx, s.DB, emailAddress, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}
