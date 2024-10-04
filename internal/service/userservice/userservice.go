package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/repository/userrepository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

type UserService struct {
	*service.App
	userRepo *userrepository.Repo
}

func New(app *service.App, userRepo *userrepository.Repo) *UserService {
	return &UserService{
		App:      app,
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, record *user.User) (*user.User, error) {
	record.OrganisationID = 1
	record.IsActive = 1

	if err := record.SetCreated(ctx); err != nil {
		s.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.userRepo.Create(ctx, tx, record)
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

func (s *UserService) GetByID(ctx context.Context, record *user.User, ID int32) error {
	err := s.userRepo.GetByID(ctx, s.DB, ID, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

func (s *UserService) GetAll(ctx context.Context, r *req.Request) ([]*user.User, error) {
	records, err := s.userRepo.GetCollection(ctx, s.DB, r)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return records, nil
}

func (s *UserService) GetByEmailAddress(ctx context.Context, record *user.User, emailAddress string) error {
	err := s.userRepo.GetByEmailAddress(ctx, s.DB, emailAddress, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}
