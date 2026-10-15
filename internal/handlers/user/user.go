// Package user provides HTTP handler functions for user resource endpoints.
package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	model "github.com/doublehops/dh-go-framework/internal/model/user"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/tools"
)

// Handle holds the dependencies for the user HTTP handler.
type Handle struct {
	srv  *userservice.UserService
	base *handlers.BaseHandler
}

// New creates a new user Handle with the required dependencies.
func New(app *service.App) *Handle {
	return &Handle{
		srv: userservice.New(app),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

// Create handles POST /user requests to create a new user record.
// Route-level role gating (org_admin/master_admin) is enforced by middleware.RequireRoles in routes.go.
func (h *Handle) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var emailAddressExists bool
	var err error
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to "+tools.CurrentFunction(), nil)

	record := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	// Organisation and role come from the requesting user, never the request body — otherwise an
	// organisation admin could create users in another organisation, or grant master_admin.
	reqUser := h.base.GetRequestingUserRecord(ctx)
	record.OrganisationID = reqUser.Record.OrganisationID

	if !reqUser.IsMasterAdmin() {
		record.RoleID = ""
	}

	if ValidationErrors := record.ValidateCreate(); len(ValidationErrors) > 0 {
		errs := req.GetValidateErrResp(ValidationErrors, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	if emailAddressExists, err = h.srv.EmailAddressAlreadyExists(ctx, record.EmailAddress); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp(req.ErrProcessingRequest.Error()))

		return
	}

	if emailAddressExists {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GeneralErrResp("email address already exists", http.StatusBadRequest))

		return
	}

	ur, err := h.srv.NewRecord(record).Create(ctx)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp(req.ErrProcessingRequest.Error()))

		return
	}

	response, err := h.GetResponse(ctx, ur.Record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

// GetCurrentUser handles GET /me. Unlike GetByID, this carries no role restriction at the route
// level (see routes.go) — any authenticated user can fetch their own record, since it's the only
// way the frontend has to learn its own role (GET /user/:id is org_admin/master_admin only).
func (h *Handle) GetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to GetCurrentUser", nil)

	reqUser := h.base.GetRequestingUserRecord(ctx)

	response, err := h.GetResponse(ctx, reqUser.Record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

// GetResponse builds a ResponseUser from a full User record, omitting sensitive fields.
func (h *Handle) GetResponse(ctx context.Context, record *model.User) (*model.ResponseUser, error) {
	userResponse := &model.ResponseUser{}
	err := copier.Copy(&userResponse, &record)
	if err != nil {
		h.base.Log.Error(ctx, "error building response object", nil)

		return nil, errors.New("error building response object")
	}

	return userResponse, nil
}

// GetCollectionResponse builds a slice of ResponseUser from a collection of User records.
func (h *Handle) GetCollectionResponse(ctx context.Context, records []*model.User) ([]*model.ResponseUser, error) {
	response := []*model.ResponseUser{}
	for _, record := range records {
		userResponse := &model.ResponseUser{}
		err := copier.Copy(&userResponse, record)
		if err != nil {
			h.base.Log.Error(ctx, "error building response object", nil)

			return nil, errors.New("error building response object")
		}

		response = append(response, userResponse)
	}

	return response, nil
}

// UpdateByID handles PATCH /user/:id requests to update an existing user record.
// Route-level role gating (org_admin/master_admin) is enforced by middleware.RequireRoles in routes.go.
func (h *Handle) UpdateByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to UpdateUser", nil)

	payload := &model.User{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	reqUser := h.base.GetRequestingUserRecord(ctx)

	id := ps.ByName("id")

	rh, err := h.srv.GetByID(ctx, id)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if rh.Record.ID == "" {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !reqUser.HasPermissionToAccessUser(rh.Record) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	if errorMsgs := payload.ValidateUpdate(); len(errorMsgs) > 0 {
		errs := req.GetValidateErrResp(errorMsgs, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	if err := rh.Update(ctx, payload); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, "Unable to process request")

		return
	}

	response, err := h.GetResponse(ctx, rh.Record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

// DeleteByID handles DELETE /user/:id requests to soft-delete a user record.
// Route-level role gating (org_admin/master_admin) is enforced by middleware.RequireRoles in routes.go.
func (h *Handle) DeleteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to DELETE user", nil)

	reqUser := h.base.GetRequestingUserRecord(ctx)

	id := ps.ByName("id")

	rh, err := h.srv.GetByID(ctx, id)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if rh.Record.ID == "" {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !reqUser.HasPermissionToAccessUser(rh.Record) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	if err := rh.Delete(ctx); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusNoContent, nil)
}

// GetByID handles GET /user/:id requests to retrieve a single user record.
// Route-level role gating (org_admin/master_admin) is enforced by middleware.RequireRoles in routes.go.
func (h *Handle) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to Get user", nil)

	reqUser := h.base.GetRequestingUserRecord(ctx)

	id := ps.ByName("id")

	rh, err := h.srv.GetByID(ctx, id)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if rh.Record.ID == "" {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !reqUser.HasPermissionToAccessUser(rh.Record) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	response, err := h.GetResponse(ctx, rh.Record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

func filterRules() []req.FilterRule {
	return []req.FilterRule{
		{
			Field: "deletedAt",
			Type:  req.FilterIsNull,
		},
		{
			Field: "name",
			Type:  req.FilterLike,
		},
	}
}

// getSortableFields will return a list of fields that a collection of records can be sorted by. This is necessary because
// not all fields should this be available to, and it will prevent SQL injection.
func getSortableFields() []string {
	return []string{
		"id",
		"name",
		"createdAt",
		"updatedAt",
	}
}

// GetAll will retrieve all records for output.
// Route-level role gating (org_admin/master_admin) is enforced by middleware.RequireRoles in routes.go.
func (h *Handle) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to Get user collection", nil)

	reqUser := h.base.GetRequestingUserRecord(ctx)

	p := req.GetRequestParams(r, filterRules(), getSortableFields())

	if !reqUser.IsMasterAdmin() {
		p.Filters = append(p.Filters, req.FilterRule{Field: "organisation_id", Type: req.FilterEquals, Value: reqUser.Record.OrganisationID})
	}

	records, err := h.srv.GetAll(ctx, p)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, "Unable to process request")

		return
	}

	response, err := h.GetCollectionResponse(ctx, records)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetListResp(response, p))
}
