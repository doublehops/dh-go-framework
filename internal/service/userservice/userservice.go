// Package userservice provides business logic for user management.
package userservice

import (
	"context"
	"errors"

	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/repository/userrepository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"
)

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

// NewRecord wraps a bare user model in a UserRecord for subsequent mutations.
func (s *UserService) NewRecord(u *user.User) *UserRecord {
	return &UserRecord{Record: u, svc: s}
}

// GetByID retrieves a user record by its ID.
func (s *UserService) GetByID(ctx context.Context, id string) (*UserRecord, error) {
	record := &user.User{}
	if err := s.userRepo.GetByID(ctx, s.DB, id, record); err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return &UserRecord{Record: record, svc: s}, nil
}

// GetAll retrieves a paginated collection of user records.
func (s *UserService) GetAll(ctx context.Context, r *req.Request) ([]*user.User, error) {
	records, err := s.userRepo.GetCollection(ctx, s.DB, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return records, nil
}

// GetByEmailAddress retrieves a user record by email address.
func (s *UserService) GetByEmailAddress(ctx context.Context, emailAddress string) (*UserRecord, error) {
	record := &user.User{}
	if err := s.userRepo.GetByEmailAddress(ctx, s.DB, emailAddress, record); err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return &UserRecord{Record: record, svc: s}, nil
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

// DeleteSessionByToken looks up a session by token and soft-deletes it.
func (s *UserService) DeleteSessionByToken(ctx context.Context, token string) error {
	us := usersessionservice.New(s.App)

	record := &usersession.UserSession{}
	if err := us.GetByToken(ctx, record, token); err != nil || record.ID == "" {
		return errors.New(service.UnableToRetrieveRecord)
	}

	return us.DeleteByID(ctx, record)
}
