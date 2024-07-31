package model

import (
	"context"
	"fmt"
	"time"

	"github.com/doublehops/dh-go-framework/internal/app"
)

type Model interface {
	GetUserID() int32
	SetCreated(context.Context) error
}

type BaseModel struct {
	ID        int32      `json:"id" db:"id"`
	UserID    int32      `json:"userId" db:"user_id"`
	CreatedBy int32      `json:"createdBy" db:"created_by"`
	UpdatedBy int32      `json:"updatedBy" db:"updated_by"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
}

// Deprecated - remove
func (bm *BaseModel) GetUserID() int32 {
	return bm.UserID
}

func (bm *BaseModel) SetCreated(ctx context.Context) error {
	userID := ctx.Value(app.UserIDKey)
	if userID != nil {
		uID, ok := userID.(int)
		if !ok {
			return fmt.Errorf("unable to convert userID to int")
		}

		bm.CreatedBy = int32(uID)
		bm.UpdatedBy = int32(uID)
		bm.UserID = int32(uID)
	}

	t := time.Now()

	bm.CreatedAt = &t
	bm.UpdatedAt = &t

	return nil
}

func (bm *BaseModel) SetUpdated(ctx context.Context) {
	userID := bm.getRequestUserID(ctx)
	if userID > 0 {
		bm.UpdatedBy = userID
	}

	t := time.Now()

	bm.UpdatedAt = &t
}

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
