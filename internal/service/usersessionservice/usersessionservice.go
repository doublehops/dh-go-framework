// Package usersessionservice provides business logic for user session management.
package usersessionservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/repository/usersessionrepository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// sessionTimeout is a sliding inactivity timeout — refreshed on every authenticated request (see
// SetLastRequestNow) — while sessionMaxAge is a hard cap on total session lifetime regardless of
// activity, so an actively-used (or leaked-and-actively-used) token can't stay valid forever.
const (
	sessionTimeout = 24 * time.Hour
	sessionMaxAge  = 7 * 24 * time.Hour
)

// UserSessionService handles business logic for user session management.
type UserSessionService struct {
	*service.App
	sessionRepo SessionRepository
}

// New creates a new UserSessionService, wiring its own repository internally.
func New(app *service.App) *UserSessionService {
	return &UserSessionService{
		App:         app,
		sessionRepo: usersessionrepository.New(app.Log),
	}
}

// Create generates a new session token and persists a session record for the given user.
func (s *UserSessionService) Create(ctx context.Context, user *user.User) (*usersession.UserSession, error) {
	token, err := generateToken(32)
	if err != nil {
		s.Log.Error(ctx, "unable to generate session token", nil)

		return nil, fmt.Errorf("unable to generate session token")
	}

	now := time.Now()
	record := &usersession.UserSession{
		UserID:      user.ID,
		Token:       token,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		LastRequest: &now,
	}

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err = s.sessionRepo.Create(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToSaveRecord+" "+err.Error(), nil)

		return record, req.ErrCouldNotSaveRecord
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	r := &usersession.UserSession{}
	err = s.sessionRepo.GetByID(ctx, s.DB, record.ID, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return r, errors.New(service.UnableToRetrieveRecord)
	}

	return r, nil
}

// UpdateLastRequest updates the last_request timestamp on a session record.
func (s *UserSessionService) UpdateLastRequest(ctx context.Context, record *usersession.UserSession) error {
	now := time.Now()
	record.LastRequest = &now
	record.UpdatedAt = &now

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.sessionRepo.UpdateLastRequest(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToUpdateRecord+" "+err.Error(), nil)

		return errors.New(service.UnableToUpdateRecord)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

// DeleteByID soft-deletes a session record by setting its deleted_at timestamp.
func (s *UserSessionService) DeleteByID(ctx context.Context, record *usersession.UserSession) error {
	now := time.Now()
	record.UpdatedAt = &now
	record.DeletedAt = &now

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.sessionRepo.Delete(ctx, tx, record)
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

// GetByID retrieves a user session record by its ID.
func (s *UserSessionService) GetByID(ctx context.Context, record *usersession.UserSession, id string) error {
	err := s.sessionRepo.GetByID(ctx, s.DB, id, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

// GetByToken retrieves a user session by its authentication token.
func (s *UserSessionService) GetByToken(ctx context.Context, record *usersession.UserSession, token string) error {
	err := s.sessionRepo.GetByToken(ctx, s.DB, token, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

// GetByUserID retrieves a user session by the associated user ID.
func (s *UserSessionService) GetByUserID(ctx context.Context, record *usersession.UserSession, userID string) error {
	err := s.sessionRepo.GetByUserID(ctx, s.DB, userID, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

// GetAll retrieves a paginated collection of user sessions.
func (s *UserSessionService) GetAll(ctx context.Context, r *req.Request) ([]*usersession.UserSession, error) {
	records, err := s.sessionRepo.GetCollection(ctx, s.DB, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return records, nil
}

func generateToken(nbytes int) (string, error) {
	b := make([]byte, nbytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// HasExpired returns true if the session's last request exceeds the session timeout duration.
func (s *UserSessionService) HasExpired(us *usersession.UserSession) bool {
	now := time.Now()

	return us.LastRequest.Add(sessionTimeout).Before(now) || us.CreatedAt.Add(sessionMaxAge).Before(now)
}

// SetLastRequestNow updates the session's last_request timestamp to the current time.
func (s *UserSessionService) SetLastRequestNow(ctx context.Context, record *usersession.UserSession) error {
	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.sessionRepo.SetLastRequestNow(ctx, s.DB, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	return nil
}
