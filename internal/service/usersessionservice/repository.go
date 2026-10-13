package usersessionservice

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

// SessionRepository defines the data access methods required by UserSessionService.
type SessionRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error
	UpdateLastRequest(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error
	Delete(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error
	GetByID(ctx context.Context, db *sqlx.DB, id string, record *usersession.UserSession) error
	GetByToken(ctx context.Context, db *sqlx.DB, token string, record *usersession.UserSession) error
	GetByUserID(ctx context.Context, db *sqlx.DB, userID string, record *usersession.UserSession) error
	GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*usersession.UserSession, error)
	SetLastRequestNow(ctx context.Context, db *sqlx.DB, record *usersession.UserSession) error
}
