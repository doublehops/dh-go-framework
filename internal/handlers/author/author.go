// Package author provides HTTP handler functions for author resource endpoints.
package author

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/authorservice"
)

// Handle holds the dependencies for the author HTTP handler.
type Handle struct {
	srv  *authorservice.AuthorService
	base *handlers.BaseHandler
}

// New creates a new author Handle with the required dependencies.
func New(app *service.App) *Handle {
	return &Handle{
		srv: authorservice.New(app),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

// Create handles POST /author requests to create a new author record.
func (h *Handle) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to CreateAuthor", nil)

	record := &author.Author{}
	if err := json.NewDecoder(r.Body).Decode(record); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	record.Normalise()

	if errs := record.Validate(); len(errs) > 0 {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GetValidateErrResp(errs, req.ErrValidation.Error()))

		return
	}

	a, err := h.srv.NewRecord(record).Create(ctx)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ServerErrResp(req.ErrProcessingRequest.Error()))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusCreated, req.GetSingleItemResp(a.Record))
}

// getAccessibleRecord loads the author named by the :id route param and checks the requesting user
// may access it. On failure it writes the 404/403 response itself and returns nil.
func (h *Handle) getAccessibleRecord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *authorservice.AuthorRecord {
	ctx := r.Context()

	a, err := h.srv.GetByID(ctx, ps.ByName("id"))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return nil
	}

	if !a.CanBeAccessedBy(h.base.GetRequestingUserRecord(ctx)) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return nil
	}

	return a
}

// UpdateByID handles PUT /author/:id requests to update an existing author record.
func (h *Handle) UpdateByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to UpdateAuthor", nil)

	a := h.getAccessibleRecord(w, r, ps)
	if a == nil {
		return
	}

	payload := &author.Author{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	payload.Normalise()

	if errs := payload.Validate(); len(errs) > 0 {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.GetValidateErrResp(errs, req.ErrValidation.Error()))

		return
	}

	if err := a.Update(ctx, payload); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(a.Record))
}

// DeleteByID handles DELETE /author/:id requests to soft-delete an author record.
func (h *Handle) DeleteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to DELETE author", nil)

	a := h.getAccessibleRecord(w, r, ps)
	if a == nil {
		return
	}

	if err := a.Delete(ctx); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusNoContent, nil)
}

// GetByID handles GET /author/:id requests to retrieve a single author record.
func (h *Handle) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to Get author", nil)

	a := h.getAccessibleRecord(w, r, ps)
	if a == nil {
		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(a.Record))
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

// GetAll handles GET /author requests to retrieve a paginated collection of author records.
func (h *Handle) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to Get authors", nil)

	p := req.GetRequestParams(r, filterRules(), getSortableFields())

	authors, err := h.srv.GetAll(ctx, p)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetListResp(authors, p))
}
