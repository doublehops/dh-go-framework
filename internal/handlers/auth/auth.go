// Package auth provides HTTP handler functions for authentication endpoints.
package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	model "github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// TokenValidPeriod is the duration after which an authentication session token expires.
const (
	TokenValidPeriod = 24 * time.Hour
)

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

	user := &model.User{}
	err = h.srv.GetByEmailAddress(ctx, user, record.EmailAddress)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, err)

		return
	}

	if passwordValid := h.srv.CheckPasswordHash(record.Password, user.Password); !passwordValid {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.ServerErrResp(req.ErrBadUsernameOrPassword.Error()))

		return
	}

	session, err := h.srv.CreateUserSession(ctx, user)
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

// func (h *Handle) GetResponse(ctx context.Context, record *model.User) (*model.ResponseUser, error) {
// 	userResponse := &model.ResponseUser{}
// 	err := copier.Copy(&userResponse, &record)
// 	if err != nil {
// 		h.base.Log.Error(ctx, "error building response object", nil)
//
// 		return nil, errors.New("error building response object")
// 	}
//
// 	return userResponse, nil
// }

// GetResponse builds the successful login response payload from the user session record.
func (h *Handle) GetResponse(_ context.Context, record *usersession.UserSession) (*model.SuccessfulLoginResponse, error) {
	response := &model.SuccessfulLoginResponse{
		BearerToken: record.Token,
	}

	return response, nil
}
