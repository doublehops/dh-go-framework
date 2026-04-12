package usersession

import (
	"time"
)

type UserSession struct {
	ID          int32      `json:"id" db:"id"`
	UserID      int32      `json:"userId" db:"user_id"`
	Token       string     `json:"token" db:"token"`
	LastRequest *time.Time `json:"lastRequest" db:"last_request"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
}
