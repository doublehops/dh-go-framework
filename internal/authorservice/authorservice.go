package authorservice

import (
	"context"

	"github.com/doublehops/dh-go-framework/internal/logga"

	"github.com/doublehops/dh-go-framework/internal/model"

	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

type AuthorService struct {
	*App
	authorRepo *repositoryauthor.Author
}

func New(app *App, authorRepo *repositoryauthor.Author) *AuthorService {
	return &AuthorService{
		App:        app,
		authorRepo: authorRepo,
	}
}

func (s AuthorService) Create(ctx context.Context, author *model.Author) (*model.Author, error) {
	if err := author.SetCreated(ctx); err != nil {
		s.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Create(ctx, tx, author)
	if err != nil {
		s.Log.Error(ctx, "unable to save new record. "+err.Error(), nil)

		return author, req.ErrCouldNotSaveRecord
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	a := &model.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, author.ID, a)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) Update(ctx context.Context, author *model.Author) (*model.Author, error) {
	author.SetUpdated(ctx)

	tx := s.DB.MustBegin()
	defer tx.Rollback() // nolint: errcheck

	err := s.authorRepo.Update(ctx, tx, author)
	if err != nil {
		s.Log.Error(ctx, "unable to update record. "+err.Error(), nil)
	}

	err = tx.Commit()
	if err != nil {
		s.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	a := &model.Author{}
	err = s.authorRepo.GetByID(ctx, s.DB, author.ID, a)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return a, nil
}

func (s AuthorService) DeleteByID(ctx context.Context, author *model.Author) error {
	tx := s.DB.MustBegin()
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

func (s AuthorService) GetByID(ctx context.Context, author *model.Author, ID int32) error {
	err := s.authorRepo.GetByID(ctx, s.DB, ID, author)
	if err != nil {
		s.Log.Error(ctx, "unable to retrieve record. "+err.Error(), nil)
	}

	return nil
}

func (s AuthorService) GetAll(ctx context.Context, p *req.Request) ([]*model.Author, error) {
	authors, err := s.authorRepo.GetCollection(ctx, s.DB, p)
	if err != nil {
		s.Log.Error(ctx, "unable to update new record. "+err.Error(), nil)
	}

	return authors, nil
}
