package authorservice

import (
	"context"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// const (
// 	unableToRetrieveRecord    = "unable to retrieve record"
// 	unableToCommitTransaction = "unable to commit transaction"
// )

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
	if err := record.SetCreated(ctx); err != nil {
		s.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Create(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToSaveRecord+" "+err.Error(), nil)

		return record, req.ErrCouldNotSaveRecord
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	a := &author.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, record.ID, a)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) Update(ctx context.Context, record *author.Author) (*author.Author, error) {
	record.SetUpdated(ctx)

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Update(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, "unable to update record. "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	a := &author.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, record.ID, a)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) DeleteByID(ctx context.Context, record *author.Author) error {
	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	record.SetDeleted(ctx)

	err := s.authorRepo.Delete(ctx, tx, record)
	if err != nil {
		s.Log.Error(ctx, "unable to delete record. "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}

func (s AuthorService) GetByID(ctx context.Context, record *author.Author, ID int32) error {
	err := s.authorRepo.GetByID(ctx, s.DB, ID, record)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return nil
}

func (s AuthorService) GetAll(ctx context.Context, p *req.Request) ([]*author.Author, error) {
	records, err := s.authorRepo.GetCollection(ctx, s.DB, p)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)
	}

	return records, nil
}
