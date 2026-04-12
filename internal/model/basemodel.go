// Package model provides the base model types and interfaces shared by all resource models.
package model

import (
	"context"
	"fmt"
	"time"

	"github.com/doublehops/dh-go-framework/internal/app"
)

// Model is the interface that all resource models must implement for permission and audit support.
type Model interface {
	GetUserID() int32
	SetCreated(context.Context) error
}

// BaseModel is embedded in all resource models and provides standard audit fields.
type BaseModel struct {
	ID        int32      `json:"id" db:"id"`
	UserID    int32      `json:"userId" db:"user_id"`
	CreatedBy int32      `json:"createdBy" db:"created_by"`
	UpdatedBy int32      `json:"updatedBy" db:"updated_by"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

// GetUserID returns the UserID field. Deprecated: to be removed.
func (bm *BaseModel) GetUserID() int32 {
	return bm.UserID
}

// SetCreated initialises the audit fields from the request context when creating a new record.
func (bm *BaseModel) SetCreated(ctx context.Context) error {
	userID := ctx.Value(app.UserIDKey)
	if userID != nil {
		uID, ok := userID.(int32)
		if !ok {
			return fmt.Errorf("unable to convert userID to int")
		}

		bm.CreatedBy = uID
		bm.UpdatedBy = uID
		bm.UserID = uID
	}

	t := time.Now()

	bm.CreatedAt = &t
	bm.UpdatedAt = &t

	return nil
}

// SetUpdated updates the UpdatedBy and UpdatedAt fields from context and current time.
func (bm *BaseModel) SetUpdated(ctx context.Context) {
	userID := bm.getRequestUserID(ctx)
	if userID > 0 {
		bm.UpdatedBy = userID
	}

	t := time.Now()

	bm.UpdatedAt = &t
}

// SetDeleted sets the DeletedAt and UpdatedAt timestamps to mark the record as soft-deleted.
func (bm *BaseModel) SetDeleted(ctx context.Context) {
	userID := bm.getRequestUserID(ctx)
	if userID > 0 {
		bm.UpdatedBy = userID
	}

	t := time.Now()

	bm.UpdatedAt = &t
	bm.DeletedAt = &t
}

// getRequestUserID will retrieve userID from context.
func (bm *BaseModel) getRequestUserID(ctx context.Context) int32 {
	val := ctx.Value(app.UserIDKey)
	var intValue int32
	var ok bool

	if intValue, ok = val.(int32); !ok {
		intValue = 0
	}

	return intValue
}
