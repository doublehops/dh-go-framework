// Package middleware provides HTTP middleware functions for authentication and request processing.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	apppackage "github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"
)

// RequireRoles returns middleware that checks the authenticated UserRecord
// (already placed in context by AuthMiddleware) holds at least one of the
// given global user_role IDs. It performs no DB lookups of its own — it only
// reads the UserRecord AuthMiddleware already fetched — so it must always be
// composed inside authMW, never used standalone. It does not, and cannot,
// check channel-scoped roles (quality_control/publisher on user_channel);
// those require a record-level check performed inside the handler instead.
func RequireRoles(roleIDs ...string) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			val := r.Context().Value(apppackage.UserRecordKey)
			reqUser, _ := val.(*userservice.UserRecord)

			if reqUser == nil || !reqUser.HasAnyRole(roleIDs...) {
				forbiddenJSON(w)

				return
			}

			next(w, r, ps)
		}
	}
}

// AuthMiddleware authenticates requests via a Bearer session token.
func AuthMiddleware(app *service.App, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.Contains(authHeader, "Bearer") {
			unauthorizedJSON(w)

			return
		}

		token, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			unauthorizedJSON(w)

			return
		}

		session := &usersession.UserSession{}

		ss := usersessionservice.New(app)
		if err := ss.GetByToken(r.Context(), session, token); err != nil {
			unauthorizedJSON(w)

			return
		}

		if ss.HasExpired(session) {
			unauthorizedJSON(w)

			return
		}

		if err := ss.SetLastRequestNow(r.Context(), session); err != nil {
			app.Log.Error(r.Context(), "failed to update last_request for session", nil)
		}

		// A deactivated (or soft-deleted, which GetByID won't find) user loses access immediately,
		// not just at their next login.
		userRecord, err := userservice.New(app).GetByID(r.Context(), session.UserID)
		if err != nil || !userRecord.IsActive() {
			unauthorizedJSON(w)

			return
		}

		ctx := context.WithValue(r.Context(), apppackage.UserIDKey, session.UserID)
		ctx = context.WithValue(ctx, apppackage.UserRecordKey, userRecord)

		next(w, r.WithContext(ctx), ps)
	}
}

func unauthorizedJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(req.GeneralErrResp(req.ErrNotAuthorised.Error(), http.StatusUnauthorized))
}

func forbiddenJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(req.GeneralErrResp(req.ErrNotAuthorised.Error(), http.StatusForbidden))
}
