// Package health provides the HTTP handler for the health check endpoint.
package health

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/handlers"
	"github.com/doublehops/dh-go-framework/internal/repository/repositoryauthor"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/authorservice"
)

// Response is the JSON response payload for the health check endpoint.
type Response struct {
	Status string `json:"status"`
}

// Handle holds dependencies for the health check handler.
type Handle struct {
	repo *repositoryauthor.Author
	srv  *authorservice.AuthorService
	base *handlers.BaseHandler
}

// New creates a new health check Handle.
func New(app *service.App) *Handle {
	ar := repositoryauthor.New(app.Log)

	return &Handle{
		repo: ar,
		srv:  authorservice.New(app, ar),
		base: &handlers.BaseHandler{
			Log: app.Log,
		},
	}
}

// Check handles GET /health requests and returns a 200 OK status response.
func (h *Handle) Check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	h.srv.Log.Info(ctx, "Request made to health check", nil)

	s := Response{
		Status: "OK",
	}

	h.base.WriteJSON(ctx, w, http.StatusOK, req.GetSingleItemResp(s))
}
