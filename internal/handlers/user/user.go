// Package user provides HTTP handler functions for user resource endpoints.
package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		srv:  userservice.New(app),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

// Create handles POST /user requests to create a new user record.
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

	if errors := record.ValidateCreate(); len(errors) > 0 {
		errs := req.GetValidateErrResp(errors, req.ErrValidation.Error())
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

	a, err := h.srv.Create(ctx, record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp(req.ErrProcessingRequest.Error()))

		return
	}

	response, err := h.GetResponse(ctx, a)
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
func (h *Handle) UpdateByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	// userID := h.base.GetUser(ctx)
	h.base.Log.Info(ctx, "Request made to UpdateUser", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	record := &model.User{}
	err = h.srv.GetByID(ctx, record, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, err)

		return
	}

	if record.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	// Uncomment to check authorization.
	// if !h.srv.HasPermission(userID, record) {
	//	 h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())
	//
	//	 return
	// }

	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	if errors := record.ValidateUpdate(); len(errors) > 0 {
		errs := req.GetValidateErrResp(errors, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	record, err = h.srv.Update(ctx, record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, "Unable to process request")

		return
	}

	response, err := h.GetResponse(ctx, record)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp("error building response object"))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(response))
}

// DeleteByID handles DELETE /user/:id requests to soft-delete a user record.
func (h *Handle) DeleteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	// userID := h.base.GetUser(ctx)
	h.base.Log.Info(ctx, "Request made to DELETE user", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	record := &model.User{}
	err = h.srv.GetByID(ctx, record, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if record.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	// Uncomment to check authorization.
	// if !h.srv.HasPermission(userID, record) {
	//	 h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())
	//
	//	 return
	// }

	if err = h.srv.DeleteByID(ctx, record); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusNoContent, nil)
}

// GetByID handles GET /user/:id requests to retrieve a single user record.
func (h *Handle) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	// userID := h.base.GetUser(ctx)
	h.base.Log.Info(ctx, "Request made to Get user", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	record := &model.User{}
	err = h.srv.GetByID(ctx, record, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if record.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	// Uncomment to check authorization.
	// if !h.srv.HasPermission(userID, record) {
	//	 h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())
	//
	//	 return
	// }

	response, err := h.GetResponse(ctx, record)
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
func (h *Handle) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.base.Log.Info(ctx, "Request made to Get user", nil)

	p := req.GetRequestParams(r, filterRules(), getSortableFields())

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
