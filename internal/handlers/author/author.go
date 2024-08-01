package author

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/doublehops/dh-go-framework/internal/model"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

type Handle struct {
	repo *repositoryauthor.Author
	srv  *service.AuthorService
	base *handlers.BaseHandler
}

func New(app *service.App) *Handle {
	ar := repositoryauthor.New(app.Log)

	return &Handle{
		repo: ar,
		srv:  service.New(app, ar),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

func (h *Handle) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to CreateAuthor", nil)

	author := &model.Author{}
	if err := json.NewDecoder(r.Body).Decode(author); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	if errors := author.Validate(); len(errors) > 0 {
		errs := req.GetValidateErrResp(errors, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	a, err := h.srv.Create(ctx, author)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.GeneralErrResp(req.ErrProcessingRequest.Error()))

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(a))
}

func (h *Handle) UpdateByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	userID := h.base.GetUser(ctx)
	h.srv.Log.Info(ctx, "Request made to UpdateAuthor", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	author := &model.Author{}
	err = h.srv.GetByID(ctx, author, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, err)

		return
	}

	if author.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !h.srv.HasPermission(userID, author) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	if err := json.NewDecoder(r.Body).Decode(author); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, req.UnableToParseResp())

		return
	}

	if errors := author.Validate(); len(errors) > 0 {
		errs := req.GetValidateErrResp(errors, req.ErrValidation.Error())
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, errs)

		return
	}

	a, err := h.srv.Update(ctx, author)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, "Unable to process request")

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(a))
}

func (h *Handle) DeleteByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	userID := h.base.GetUser(ctx)
	h.srv.Log.Info(ctx, "Request made to DELETE author", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	author := &model.Author{}
	err = h.srv.GetByID(ctx, author, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if author.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !h.srv.HasPermission(userID, author) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	if err = h.srv.DeleteByID(ctx, author); err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, req.ErrorProcessingRequestResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusNoContent, nil)
}

func (h *Handle) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	userID := h.base.GetUser(ctx)
	h.srv.Log.Info(ctx, "Request made to Get author", nil)

	ID := ps.ByName("id")
	i, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusBadRequest, "ID is not a valid value")

		return
	}

	author := &model.Author{}
	err = h.srv.GetByID(ctx, author, int32(i))
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, "Unable to find record")

		return
	}

	if author.ID == 0 {
		h.base.WriteJSON(ctx, w, http.StatusNotFound, req.GetNotFoundResp())

		return
	}

	if !h.srv.HasPermission(userID, author) {
		h.base.WriteJSON(ctx, w, http.StatusForbidden, req.GetNotAuthorisedResp())

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(author))
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

func (h *Handle) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to Get authors", nil)

	p := req.GetRequestParams(r, filterRules(), getSortableFields())

	authors, err := h.srv.GetAll(ctx, p)
	if err != nil {
		h.base.WriteJSON(ctx, w, http.StatusInternalServerError, "Unable to process request")

		return
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetListResp(authors, p))
}
