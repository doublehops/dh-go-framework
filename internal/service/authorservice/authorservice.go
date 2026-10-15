// Package authorservice provides business logic for author management.
package authorservice

import (
	"context"

	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// AuthorService handles business logic for author management. Collection-level operations live
// here; operations on a single record live on AuthorRecord (see authorrecord.go).
type AuthorService struct {
	*service.App
	authorRepo AuthorRepository
}

// New creates a new AuthorService, wiring its own repository internally.
func New(app *service.App) *AuthorService {
	return &AuthorService{
		App:        app,
		authorRepo: repositoryauthor.New(app.Log),
	}
}

// NewRecord wraps a bare author model in an AuthorRecord.
func (s *AuthorService) NewRecord(a *author.Author) *AuthorRecord {
	return &AuthorRecord{Record: a, svc: s}
}

// GetByID retrieves an author record by its ID. An error is returned when it doesn't exist.
func (s *AuthorService) GetByID(ctx context.Context, id string) (*AuthorRecord, error) {
	record := &author.Author{}
	if err := s.authorRepo.GetByID(ctx, s.DB, id, record); err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return s.NewRecord(record), nil
}

// GetAll retrieves a paginated collection of author records.
func (s *AuthorService) GetAll(ctx context.Context, p *req.Request) ([]*author.Author, error) {
	records, err := s.authorRepo.GetCollection(ctx, s.DB, p)
	if err != nil {
		s.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return nil, err
	}

	return records, nil
}
