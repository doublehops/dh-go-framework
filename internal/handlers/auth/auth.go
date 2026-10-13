// Package auth provides HTTP handler functions for authentication endpoints.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	model "github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// dummyPasswordHash is compared against on a login attempt for an email address that doesn't
// exist, so the request takes roughly the same time as a real bcrypt comparison would — without
// this, a nonexistent-email request returns near-instantly (just a failed DB lookup) while a
// real-but-wrong-password request takes bcrypt's ~100ms, letting an attacker distinguish "no such
// account" from "wrong password" purely by response time despite both returning the same error
// message. The plaintext this hash was generated from is never used for a real account.
const dummyPasswordHash = "$2a$10$0FHn7e7cMTRhvEZlpOk5bulC6Qhju8HWirXdq2yebzdig2ycf9Mu6"

// Handle holds the dependencies for the auth handler.
type Handle struct {
	srv  *userservice.UserService
	base *handlers.BaseHandler
}

// New creates a new auth Handle with the required dependencies.
func New(app *service.App) *Handle {
	return &Handle{
		srv: userservice.New(app),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

// Login handles POST /auth/login requests, authenticating the user and returning a session token.
func (h *Handle) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to "+tools.CurrentFunction(), nil)

	record := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	if errs := record.ValidateLogin(); len(errs) > 0 {
		errs := req.GetValidateErrResp(errs, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	userRecord, err := h.srv.GetByEmailAddress(ctx, record.EmailAddress)
	if err != nil {
		// Still runs a bcrypt comparison (against a fixed dummy hash, its result discarded) so
		// this response takes about as long as the real-user-wrong-password path below — see
		// dummyPasswordHash's comment.
		h.srv.CheckPasswordHash(record.Password, dummyPasswordHash)
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GeneralErrResp(req.ErrBadUsernameOrPassword.Error(), http.StatusBadRequest))

		return
	}

	if !userRecord.CheckPasswordHash(record.Password) {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GeneralErrResp(req.ErrBadUsernameOrPassword.Error(), http.StatusBadRequest))

		return
	}

	if !userRecord.IsActive() {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GeneralErrResp(req.ErrUserNotActive.Error(), http.StatusBadRequest))

		return
	}

	session, err := userRecord.CreateSession(ctx)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("unable to create session"))

		return
	}

	response, err := h.GetResponse(ctx, session)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

// Logout handles POST /auth/logout requests, invalidating the caller's session token.
func (h *Handle) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to "+tools.CurrentFunction(), nil)

	token, _ := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")

	if err := h.srv.DeleteSessionByToken(ctx, token); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("unable to logout"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusNoContent, nil)
}

// GetResponse builds the successful login response payload from the user session record.
func (h *Handle) GetResponse(_ context.Context, record *usersession.UserSession) (*model.SuccessfulLoginResponse, error) {
	response := &model.SuccessfulLoginResponse{
		BearerToken: record.Token,
	}

	return response, nil
}
