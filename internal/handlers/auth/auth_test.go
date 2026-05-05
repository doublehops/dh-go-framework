package auth

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

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	usermodel "github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/testtools"
)

const testPassword = "pass123"

var (
	cfg       *config.Config
	testEmail string
)

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.New("./config_test.json")
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	stopServer, err := testtools.StartTestServer()
	if err != nil {
		log.Fatalf("failed to start test server: %v", err)
	}

	if err = testtools.RunMigrations(cfg); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	testEmail = fmt.Sprintf("authtest-%d@example.com", time.Now().UnixNano())

	if _, _, err = testtools.CreateTestUser(cfg, testEmail, testPassword); err != nil {
		log.Fatalf("failed to create test user: %v", err)
	}

	code := m.Run()
	stopServer()
	os.Exit(code)
}

// login is a helper that POSTs to /v1/auth/login and returns the bearer token on success.
func login(r httprequest.Requester, email, password string) (string, string, error) {
	payload := usermodel.User{
		EmailAddress: email,
		Password:     password,
	}

	statusCode, body, err := r.MakeRequest(context.Background(), http.MethodPost, "v1/auth/login", nil, payload)
	if err != nil {
		return statusCode, "", err
	}

	resp := request.SingleItemResp{Data: &usermodel.SuccessfulLoginResponse{}}
	if err = json.Unmarshal(body, &resp); err != nil {
		return statusCode, "", err
	}

	loginResp, ok := resp.Data.(*usermodel.SuccessfulLoginResponse)
	if !ok || loginResp.BearerToken == "" {
		return statusCode, "", nil
	}

	return statusCode, "Bearer " + loginResp.BearerToken, nil
}

func TestLogin(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)

	// Valid credentials.
	statusCode, token, err := login(r, testEmail, testPassword)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))
	assert.NotEmpty(t, token)

	// Wrong password.
	statusCode, _, err = login(r, testEmail, "wrongpassword")
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))

	// Non-existent email.
	statusCode, _, err = login(r, "nobody@example.com", testPassword)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))

	// Missing email (validation failure).
	statusCode, _, err = login(r, "", testPassword)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))

	// Missing password (validation failure).
	statusCode, _, err = login(r, testEmail, "")
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))
}

func TestLogout(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	// Obtain a fresh token for this test.
	statusCode, token, err := login(r, testEmail, testPassword)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	authHeader := map[string]string{"Authorization": token}

	// Without auth — middleware should reject.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodPost, "v1/auth/logout", nil, nil)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// With valid token — should succeed.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodPost, "v1/auth/logout", nil, nil, authHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	// Same token again — session deleted, middleware should reject.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodPost, "v1/auth/logout", nil, nil, authHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))
}
