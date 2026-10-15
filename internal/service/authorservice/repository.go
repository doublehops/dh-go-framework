package authorservice

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/model/author"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

// AuthorRepository defines the data access methods required by AuthorService.
type AuthorRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, record *author.Author) error
	Update(ctx context.Context, tx *sqlx.Tx, record *author.Author) error
	Delete(ctx context.Context, tx *sqlx.Tx, record *author.Author) error
	GetByID(ctx context.Context, db *sqlx.DB, id string, record *author.Author) error
	GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*author.Author, error)
}
