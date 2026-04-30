package testtools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/logga"
	usermodel "github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/repository/userrepository"
	"github.com/doublehops/dh-go-framework/internal/repository/usersessionrepository"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/service/usersessionservice"
)

// TestCredentials holds the email/password for a test user created during setup.
type TestCredentials struct {
	Email    string
	Password string
}

// CreateTestUser creates a user and session directly via the service layer, bypassing HTTP.
// Returns the Authorization header value ("Bearer <token>") ready to pass to MakeRequest.
func CreateTestUser(cfg *config.Config, email, password string) (string, error) {
	logCfg := &config.Logging{
		Writer:       "stdout",
		LogLevel:     "DEBUG",
		OutputFormat: "JSON",
	}

	l, err := logga.New(logCfg)
	if err != nil {
		return "", fmt.Errorf("CreateTestUser: failed to create logger: %w", err)
	}

	database, err := db.New(l, cfg.DB)
	if err != nil {
		return "", fmt.Errorf("CreateTestUser: failed to connect to database: %w", err)
	}

	app := &service.App{
		DB:  database,
		Log: l,
	}

	ctx := context.Background()

	userRepo := userrepository.New(l)
	userSvc := userservice.New(app, userRepo)

	record := &usermodel.User{
		Name:         "testuser",
		EmailAddress: email,
		Password:     password,
	}

	user, err := userSvc.Create(ctx, record)
	if err != nil {
		return "", fmt.Errorf("CreateTestUser: failed to create user: %w", err)
	}

	sessionRepo := usersessionrepository.New(l)
	sessionSvc := usersessionservice.New(app, sessionRepo)

	session, err := sessionSvc.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("CreateTestUser: failed to create session: %w", err)
	}

	return "Bearer " + session.Token, nil
}

// GetAuthToken logs in with the given credentials and returns the Authorization header value
// ("Bearer <token>") ready to pass directly to MakeRequest.
// Returns an error rather than calling t.Fatal so it can be used from TestMain.
func GetAuthToken(req *httprequest.Requester, creds TestCredentials) (string, error) {
	payload := usermodel.User{
		EmailAddress: creds.Email,
		Password:     creds.Password,
	}

	statusCode, body, err := req.MakeRequest(context.Background(), http.MethodPost, "v1/auth/login", nil, payload)
	if err != nil {
		return "", fmt.Errorf("GetAuthToken: request failed: %w", err)
	}

	if !strings.Contains(statusCode, fmt.Sprintf("%d", http.StatusOK)) {
		return "", fmt.Errorf("GetAuthToken: unexpected status %s, body: %s", statusCode, body)
	}

	resp := request.SingleItemResp{
		Data: &usermodel.SuccessfulLoginResponse{},
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("GetAuthToken: failed to unmarshal response: %w", err)
	}

	loginResp, ok := resp.Data.(*usermodel.SuccessfulLoginResponse)
	if !ok || loginResp.BearerToken == "" {
		return "", fmt.Errorf("GetAuthToken: empty or missing bearer token in response")
	}

	return "Bearer " + loginResp.BearerToken, nil
}

// MustGetAuthToken is a test-only convenience wrapper around GetAuthToken that calls t.Fatal on error.
// Use this within individual tests; use GetAuthToken directly in TestMain.
func MustGetAuthToken(t *testing.T, req *httprequest.Requester, creds TestCredentials) string {
	t.Helper()

	token, err := GetAuthToken(req, creds)
	if err != nil {
		t.Fatalf("MustGetAuthToken: %v", err)
	}

	return token
}
