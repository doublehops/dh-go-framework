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

func TestAuthorCreate(t *testing.T) {
	var ok bool
	var d *author.Author
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	authHeader := map[string]string{"Authorization": authToken}

	payload := author.Author{
		Name: "author1",
	}

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

func TestAuthorUpdate(t *testing.T) {
	var ok bool

	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	ctx = context.WithValue(ctx, app.UserIDKey, authedUser.ID)

	authHeader := map[string]string{"Authorization": authToken}

	ar := &author.Author{
		Name: "authorOriginal",
	}
	appObj, err := testtools.CreateApp()
	assert.NoError(t, err, "unexpected error creating app")

	srv := authorservice.New(appObj)
	record, err := srv.Create(ctx, ar)
	assert.NoError(t, err, "unexpected error creating record")

	response := request.SingleItemResp{
		Data: &author.Author{},
	}
	// Test UPDATE new record.
	payload := author.Author{
		Name: "authorUpdated",
	}

	// Without auth.
	path := fmt.Sprintf("v1/author/%d", record.ID)
	statusCode, res, err := req.MakeRequest(ctx, http.MethodPut, path, nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))
	// With auth.
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPut, path, nil, payload, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &response)
	assert.NoError(t, err, "unable to unmarshal record")
	if record, ok = response.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, record.Name)
	assert.Equal(t, record.ID, record.ID)
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *record.UpdatedAt, duration)
}

func TestAuthorDelete(t *testing.T) {
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	ctx = context.WithValue(ctx, app.UserIDKey, authedUser.ID)

	authHeader := map[string]string{"Authorization": authToken}

	ar := &author.Author{
		Name: "authorOriginal",
	}
	appObj, err := testtools.CreateApp()
	assert.NoError(t, err, "unexpected error creating app")

	srv := authorservice.New(appObj)
	record, err := srv.Create(ctx, ar)
	assert.NoError(t, err, "unexpected error creating record")

	// Without auth.
	path := fmt.Sprintf("v1/author/%d", record.ID)
	statusCode, _, err := req.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth.
	statusCode, _, err = req.MakeRequest(ctx, http.MethodDelete, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	statusCode, _, err = req.MakeRequest(ctx, http.MethodDelete, path, nil, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))
}

func TestAuthorCollection(t *testing.T) {
	type collectionResp struct {
		Data       []*author.Author `json:"data"`
		TotalCount int32            `json:"totalCount"`
	}

	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	ctx = context.WithValue(ctx, app.UserIDKey, authedUser.ID)

	authHeader := map[string]string{"Authorization": authToken}

	appObj, err := testtools.CreateApp()
	assert.NoError(t, err, "unexpected error creating app")

	srv := authorservice.New(appObj)

	ar1, err := srv.Create(ctx, &author.Author{Name: "collectionAuthor1"})
	assert.NoError(t, err, "unexpected error creating record")
	ar2, err := srv.Create(ctx, &author.Author{Name: "collectionAuthor2"})
	assert.NoError(t, err, "unexpected error creating record")

	// Without auth.
	statusCode, _, err := req.MakeRequest(ctx, http.MethodGet, "v1/author", nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With auth — sort newest first so the just-created records appear on page 1.
	params := map[string]string{"sortBy": "id", "order": "desc"}
	statusCode, res, err := req.MakeRequest(ctx, http.MethodGet, "v1/author", params, nil, authHeader)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	var response collectionResp
	err = json.Unmarshal(res, &response)
	assert.NoError(t, err, "unable to unmarshal collection response")

	assert.GreaterOrEqual(t, response.TotalCount, int32(2))

	var names []string
	for _, a := range response.Data {
		names = append(names, a.Name)
	}
	assert.Contains(t, names, ar1.Name)
	assert.Contains(t, names, ar2.Name)
}
