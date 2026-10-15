package userservice

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/model/user"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

// UserRepository defines the data access methods required by UserService.
type UserRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	Update(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	Delete(ctx context.Context, tx *sqlx.Tx, record *user.User) error
	GetByID(ctx context.Context, db *sqlx.DB, id string, record *user.User) error
	GetByEmailAddress(ctx context.Context, db *sqlx.DB, emailAddress string, record *user.User) error
	GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*user.User, error)
}
