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

	"github.com/doublehops/dh-go-framework/internal/app"
	"github.com/doublehops/dh-go-framework/internal/service/authorservice"
	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/testtools"
)

var (
	cfg        *config.Config
	authToken  string // Authorization header value set once in TestMain
	authedUser *user.User
)

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.New("./config_test.json")
	if err != nil {
		log.Printf("error starting main. %s", err.Error())
		os.Exit(1)
	}

	stopServer, err := testtools.StartTestServer()
	if err != nil {
		log.Fatalf("failed to start test server: %v", err)
	}

	if err = testtools.RunMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Use a unique email per test run to avoid conflicts with existing records.
	email := fmt.Sprintf("testuser-%d@example.com", time.Now().UnixNano())

	authedUser, authToken, err = testtools.CreateTestUser(cfg, email, "pass123")
	if err != nil {
		log.Fatalf("failed to create test user: %v", err)
	}

	code := m.Run()
	stopServer()
	os.Exit(code)
}

//nolint:funlen
func TestAuthorCRUD_TODO_SPLIT(t *testing.T) {
	var ok bool
	var d *author.Author
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	authHeader := map[string]string{"Authorization": authToken}

	payload := author.Author{
		Name: "author1",
	}

	// Test CREATE new record.
	// Without auth.
	statusCode, res, err := req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	record := request.SingleItemResp{
		Data: &author.Author{},
	}
	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.CreatedAt, duration)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test GET new record.
	// Without auth.
	path := fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))
	// With auth.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodGet, path, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test UPDATE new record.
	payload = author.Author{
		Name: "authorABC",
	}

	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPut, path, nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test DELETE new record.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, _, err = req.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	// Test that record has been deleted.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, _, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))

	// Test that GetAll requires authentication.
	statusCode, _, err = req.MakeRequest(ctx, http.MethodGet, "v1/author", nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// Test that GetAll succeeds with a valid token.
	statusCode, _, err = req.MakeRequest(ctx, http.MethodGet, "v1/author", nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))
}

func TestAuthorCreate(t *testing.T) {
	var ok bool
	var d *author.Author
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	authHeader := map[string]string{"Authorization": authToken}

	payload := author.Author{
		Name: "author1",
	}

	// Test CREATE new record.
	// Without auth.
	statusCode, res, err := req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	record := request.SingleItemResp{
		Data: &author.Author{},
	}
	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.CreatedAt, duration)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)
}

func TestAuthorGet(t *testing.T) {
	var ok bool
	var d *author.Author
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	ctx = context.WithValue(ctx, app.UserIDKey, authedUser.ID)

	authHeader := map[string]string{"Authorization": authToken}

	ar := &author.Author{
		Name: "author1",
	}
	appObj, err := testtools.CreateApp()
	assert.NoError(t, err, "unexpected error creating app")

	srv := authorservice.New(appObj)
	record, err := srv.Create(ctx, ar)
	assert.NoError(t, err, "unexpected error creating record")
	response := request.SingleItemResp{
		Data: &author.Author{},
	}

	// Test GET new record.
	// Without auth.
	path := fmt.Sprintf("v1/author/%d", record.ID)
	statusCode, _, err := req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))
	// With auth.
	path = fmt.Sprintf("v1/author/%d", record.ID)
	statusCode, res, err := req.MakeRequest(ctx, http.MethodGet, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &response)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = response.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, ar.Name, d.Name)
	assert.Equal(t, d.ID, record.ID)
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)
}
