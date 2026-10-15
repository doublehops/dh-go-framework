// Package handlers provides base HTTP handler utilities shared across all resource handlers.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
)

// BaseHandler provides common handler utilities like JSON response writing and user retrieval.
type BaseHandler struct {
	Log *logga.Logga
}

// GetUser retrieves the authenticated user ID from the request context.
func (bh *BaseHandler) GetUser(ctx context.Context) string {
	val := ctx.Value(app.UserIDKey)
	strValue, _ := val.(string)

	return strValue
}

// GetRequestingUserRecord retrieves the authenticated UserRecord from the request context.
func (bh *BaseHandler) GetRequestingUserRecord(ctx context.Context) *userservice.UserRecord {
	val := ctx.Value(app.UserRecordKey)
	record, _ := val.(*userservice.UserRecord)

	return record
}

// WriteJSON writes a JSON-encoded response with the given status code.
func (bh *BaseHandler) WriteJSON(ctx context.Context, w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		bh.Log.Error(ctx, "unable to marshal to JSON. "+err.Error(), nil)
	}
}
