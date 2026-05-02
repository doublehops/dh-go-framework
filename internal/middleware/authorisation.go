// Package middleware provides HTTP middleware functions for authentication and request processing.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	apppackage "github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"
)

// AuthMiddleware will authenticate user by the bearer token passed in through the authorization header.
// todo - needs implementation.
func AuthMiddleware(app *service.App, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.Contains(authHeader, "Bearer") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		token, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		record := &usersession.UserSession{}

		ss := usersessionservice.New(app)
		err := ss.GetByToken(r.Context(), record, token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		if ss.HasExpired(record) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		// TODO - handle error
		_ = ss.SetLastRequestNow(r.Context(), record)

		r = r.WithContext(context.WithValue(r.Context(), apppackage.UserIDKey, record.UserID))
		next(w, r, ps)
	}
}
