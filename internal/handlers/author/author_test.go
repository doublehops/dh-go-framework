package author

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service/authorservice"
	"github.com/doublehops/dh-go-framework/internal/testtools"
)

var (
	cfg        *config.Config
	authToken  string // Authorization header value set once in TestMain
	authedUser *user.User
	// otherAuthToken belongs to a second, non-admin user who doesn't own the records created here.
	otherAuthToken string
)

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.New("./config_test.json")
	if err != nil {
		log.Printf("error starting main. %s", err.Error())
		os.Exit(1)
	}

	stopServer, err := testtools.StartTestServer(cfg)
	if err != nil {
		log.Fatalf("failed to start test server: %v", err)
	}

	// Once the server subprocess is up, any later fatal setup failure must stop it before
	// exiting — log.Fatalf calls os.Exit directly, which skips deferred cleanup.
	fatalf := func(format string, args ...any) {
		stopServer()
		log.Fatalf(format, args...)
	}

	if err = testtools.RunMigrations(cfg); err != nil {
		fatalf("failed to run migrations: %v", err)
	}

	// Use a unique email per test run to avoid conflicts with existing records.
	email := fmt.Sprintf("testuser-%d@example.com", time.Now().UnixNano())

	authedUser, authToken, err = testtools.CreateTestUser(cfg, email, "pass1234")
	if err != nil {
		fatalf("failed to create test user: %v", err)
	}

	otherEmail := fmt.Sprintf("otheruser-%d@example.com", time.Now().UnixNano())

	if _, otherAuthToken, err = testtools.CreateTestUser(cfg, otherEmail, "pass1234"); err != nil {
		fatalf("failed to create second test user: %v", err)
	}

	code := m.Run()
	stopServer()
	os.Exit(code)
}

// setup wires up the common test dependencies: HTTP requester, a context carrying
// the authenticated user ID, the Authorization header, and an authorservice instance
// for direct record creation (bypassing HTTP).
func setup(t *testing.T) (httprequest.Requester, context.Context, map[string]string, *authorservice.AuthorService) {
	t.Helper()

	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.WithValue(context.Background(), app.UserIDKey, authedUser.ID)
	authHeader := map[string]string{"Authorization": authToken}

	appObj, err := testtools.CreateApp()
	assert.NoError(t, err, "unexpected error creating app")

	return r, ctx, authHeader, authorservice.New(appObj)
}

// decodeAuthor unmarshals a single-item JSON response body into an *author.Author.
func decodeAuthor(t *testing.T, body []byte) *author.Author {
	t.Helper()

	resp := request.SingleItemResp{Data: &author.Author{}}
	assert.NoError(t, json.Unmarshal(body, &resp), "unable to unmarshal response")

	d, ok := resp.Data.(*author.Author)
	if !ok {
		t.Fatal("unable to type-assert response data to *author.Author")
	}

	return d
}

func TestAuthorCreate(t *testing.T) {
	r, ctx, authHeader, _ := setup(t)

	payload := author.Author{Name: "author1"}

	// Without auth.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, res, err := r.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	d := decodeAuthor(t, res)
	assert.Equal(t, payload.Name, d.Name)
	_, err = uuid.Parse(d.ID)
	assert.NoError(t, err, "expected ID to be a UUID")
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.CreatedAt, duration)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)
}

func TestAuthorGet(t *testing.T) {
	r, ctx, authHeader, srv := setup(t)

	record, err := srv.NewRecord(&author.Author{Name: "author1"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")

	path := "v1/author/" + record.Record.ID

	// Without auth.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, res, err := r.MakeRequest(ctx, http.MethodGet, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	d := decodeAuthor(t, res)
	assert.Equal(t, record.Record.Name, d.Name)
	assert.Equal(t, record.Record.ID, d.ID)
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)
}

func TestAuthorUpdate(t *testing.T) {
	r, ctx, authHeader, srv := setup(t)

	record, err := srv.NewRecord(&author.Author{Name: "authorOriginal"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")

	payload := author.Author{Name: "authorUpdated"}
	path := "v1/author/" + record.Record.ID

	// Without auth.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodPut, path, nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, res, err := r.MakeRequest(ctx, http.MethodPut, path, nil, payload, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	d := decodeAuthor(t, res)
	assert.Equal(t, payload.Name, d.Name)
	assert.Equal(t, record.Record.ID, d.ID)
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)
}

func TestAuthorDelete(t *testing.T) {
	r, ctx, authHeader, srv := setup(t)

	record, err := srv.NewRecord(&author.Author{Name: "authorOriginal"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")

	path := "v1/author/" + record.Record.ID

	// Without auth.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodDelete, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	// Already deleted — should now be not found.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodDelete, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))
}

func TestAuthorCollection(t *testing.T) {
	type collectionResp struct {
		Data       []*author.Author `json:"data"`
		TotalCount int32            `json:"totalCount"`
	}

	r, ctx, authHeader, srv := setup(t)

	ar1, err := srv.NewRecord(&author.Author{Name: "collectionAuthor1"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")
	ar2, err := srv.NewRecord(&author.Author{Name: "collectionAuthor2"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")

	// Without auth.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, "v1/author", nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth — sort newest first so the just-created records appear on page 1. (UUIDs are
	// random, so sorting by id no longer gives creation order.)
	params := map[string]string{"sortBy": "createdAt", "order": "desc"}
	statusCode, res, err := r.MakeRequest(ctx, http.MethodGet, "v1/author", params, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	var response collectionResp
	assert.NoError(t, json.Unmarshal(res, &response), "unable to unmarshal collection response")
	assert.GreaterOrEqual(t, response.TotalCount, int32(2))

	var names []string
	for _, a := range response.Data {
		names = append(names, a.Name)
	}
	assert.Contains(t, names, ar1.Record.Name)
	assert.Contains(t, names, ar2.Record.Name)
}

func TestAuthorNotFound(t *testing.T) {
	r, ctx, authHeader, _ := setup(t)

	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, "v1/author/"+uuid.New().String(), nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))
}

func TestAuthorForbiddenForNonOwner(t *testing.T) {
	r, ctx, _, srv := setup(t)

	record, err := srv.NewRecord(&author.Author{Name: "ownedAuthor"}).Create(ctx)
	assert.NoError(t, err, "unexpected error creating record")

	path := "v1/author/" + record.Record.ID
	otherHeader := map[string]string{"Authorization": otherAuthToken}

	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, path, nil, nil, otherHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusForbidden))

	statusCode, _, err = r.MakeRequest(ctx, http.MethodPut, path, nil, author.Author{Name: "hijacked"}, otherHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusForbidden))

	statusCode, _, err = r.MakeRequest(ctx, http.MethodDelete, path, nil, nil, otherHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusForbidden))
}

func TestAuthorCreateNormalisesBeforeValidating(t *testing.T) {
	r, ctx, authHeader, _ := setup(t)

	// All whitespace: trimmed to "" by Normalise, so it fails the required/length rule.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodPost, "v1/author", nil, author.Author{Name: "     "}, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))

	statusCode, res, err := r.MakeRequest(ctx, http.MethodPost, "v1/author", nil, author.Author{Name: "  padded  "}, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))
	assert.Equal(t, "padded", decodeAuthor(t, res).Name)
}
