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

const sessionTimeout = 24 * time.Hour

type UserSessionService struct {
	*service.App
	sessionRepo *usersessionrepository.Repo
}

func New(app *service.App, sessionRepo *usersessionrepository.Repo) *UserSessionService {
	return &UserSessionService{
		App:         app,
		sessionRepo: sessionRepo,
	}
}

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

func (s *UserSessionService) UpdateLastRequest(ctx context.Context, record *usersession.UserSession) error {
	now := time.Now()
	record.LastRequest = &now
	record.UpdatedAt = &now

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.sessionRepo.UpdateLastRequest(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToUpdateRecord+" "+err.Error(), nil)

		return fmt.Errorf(service.UnableToUpdateRecord)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

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

func (s *UserSessionService) GetByID(ctx context.Context, record *usersession.UserSession, ID int32) error {
	err := s.sessionRepo.GetByID(ctx, s.DB, ID, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

func (s *UserSessionService) GetByToken(ctx context.Context, record *usersession.UserSession, token string) error {
	err := s.sessionRepo.GetByToken(ctx, s.DB, token, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

func (s *UserSessionService) GetByUserID(ctx context.Context, record *usersession.UserSession, userID int32) error {
	err := s.sessionRepo.GetByUserID(ctx, s.DB, userID, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

func (s *UserSessionService) GetAll(ctx context.Context, r *req.Request) ([]*usersession.UserSession, error) {
	records, err := s.sessionRepo.GetCollection(ctx, s.DB, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
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

func (s *UserSessionService) HasExpired(us *usersession.UserSession) bool {
	return us.LastRequest.Add(sessionTimeout).Before(time.Now())
}

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
