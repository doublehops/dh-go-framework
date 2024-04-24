package authorservice

import (
	"context"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/service"

	"github.com/doublehops/dh-go-framework/internal/logga"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

type AuthorService struct {
	*service.App
	authorRepo *repositoryauthor.Author
}

func New(app *service.App, authorRepo *repositoryauthor.Author) *AuthorService {
	return &AuthorService{
		App:        app,
		authorRepo: authorRepo,
	}
}

func (s AuthorService) Create(ctx context.Context, record *author.Author) (*author.Author, error) {
	ctx = context.WithValue(ctx, app.UserIDKey, 1) // todo - set this in middleware.

	if err := record.SetCreated(ctx); err != nil {
		s.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx, _ := s.DB.BeginTx(ctx, nil)
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Create(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, "unable to save new record. "+err.Error(), nil)

		return record, req.ErrCouldNotSaveRecord
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	a := &author.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, record.ID, a)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) Update(ctx context.Context, record *author.Author) (*author.Author, error) {
	record.SetUpdated(ctx)

	tx, _ := s.DB.BeginTx(ctx, nil)
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Update(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, "unable to update record. "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	a := &author.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, record.ID, a)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) DeleteByID(ctx context.Context, author *author.Author) error {
	tx, _ := s.DB.BeginTx(ctx, nil)
	defer tx.Rollback() // nolint: errcheck

	author.SetDeleted(ctx)

	err := s.authorRepo.Delete(ctx, tx, author)
	if err != nil {
		s.Log.Error(ctx, "unable to delete record. "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	return nil
}

func (s AuthorService) GetByID(ctx context.Context, record *author.Author, ID int32) error {
	err := s.authorRepo.GetByID(ctx, s.DB, ID, record)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return nil
}

func (s AuthorService) GetAll(ctx context.Context, p *req.Request) ([]*author.Author, error) {
	authors, err := s.authorRepo.GetAll(ctx, s.DB, p)
	if err != nil {
		s.Log.Error(ctx, "unable to update new record. "+err.Error(), nil)
	}

	return authors, nil
}
