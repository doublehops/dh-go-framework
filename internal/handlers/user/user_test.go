package user

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
	"github.com/doublehops/dh-go-framework/internal/model/userrole"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service/userservice"
	"github.com/doublehops/dh-go-framework/internal/testtools"
)

const testPassword = "pass1234"

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.New("./config_test.json")
	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	stopServer, err := testtools.StartTestServer(cfg)
	if err != nil {
		log.Fatalf("failed to start test server: %v", err)
	}

	if err = testtools.RunMigrations(cfg); err != nil {
		stopServer()
		log.Fatalf("failed to run migrations: %v", err)
	}

	code := m.Run()
	stopServer()
	os.Exit(code)
}

// uniqueEmail returns an email address that won't collide with earlier runs against the same DB.
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano())
}

// createUser creates a user with the given role via testtools and returns it with its auth header.
func createUser(t *testing.T, prefix, roleID string) (*usermodel.User, map[string]string) {
	t.Helper()

	u, token, err := testtools.CreateTestUserWithRole(cfg, uniqueEmail(prefix), testPassword, roleID)
	assert.NoError(t, err, "unexpected error creating test user")

	return u, map[string]string{"Authorization": token}
}

func decodeUser(t *testing.T, body []byte) *usermodel.ResponseUser {
	t.Helper()

	resp := request.SingleItemResp{Data: &usermodel.ResponseUser{}}
	assert.NoError(t, json.Unmarshal(body, &resp), "unable to unmarshal response")

	u, ok := resp.Data.(*usermodel.ResponseUser)
	if !ok {
		t.Fatal("unable to type-assert response data to *usermodel.ResponseUser")
	}

	return u
}

func TestUserRoutesRequireAdminRole(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	standard, standardHeader := createUser(t, "standard", "")
	_, adminHeader := createUser(t, "orgadmin", userrole.OrganisationAdminID)

	// A standard user is refused every /user route.
	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, "v1/user", nil, nil, standardHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusForbidden))

	statusCode, _, err = r.MakeRequest(ctx, http.MethodGet, "v1/user/"+standard.ID, nil, nil, standardHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusForbidden))

	// An organisation admin can read users in their organisation.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodGet, "v1/user", nil, nil, adminHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	statusCode, res, err := r.MakeRequest(ctx, http.MethodGet, "v1/user/"+standard.ID, nil, nil, adminHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))
	assert.Equal(t, standard.EmailAddress, decodeUser(t, res).EmailAddress)
}

func TestGetCurrentUser(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, "v1/me", nil, nil)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// /me needs no admin role, and reports the default role for a user created without one.
	standard, standardHeader := createUser(t, "me", "")

	statusCode, res, err := r.MakeRequest(ctx, http.MethodGet, "v1/me", nil, nil, standardHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	me := decodeUser(t, res)
	assert.Equal(t, standard.ID, me.ID)
	assert.Equal(t, userrole.StandardUserID, me.RoleID)
	assert.NotContains(t, string(res), "password", "password hash must not be returned")
}

func TestUserCreateIgnoresRoleAndOrganisationFromBody(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	admin, adminHeader := createUser(t, "orgadmin", userrole.OrganisationAdminID)

	// An organisation admin must not be able to mint a master admin, or place a user elsewhere.
	payload := usermodel.User{
		Name:           "newuser",
		EmailAddress:   uniqueEmail("created"),
		Password:       testPassword,
		RoleID:         userrole.MasterAdminID,
		OrganisationID: "00000000-0000-0000-0000-000000000999",
	}

	statusCode, res, err := r.MakeRequest(ctx, http.MethodPost, "v1/user", nil, payload, adminHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	created := decodeUser(t, res)
	assert.Equal(t, userrole.StandardUserID, created.RoleID)
	assert.Equal(t, admin.OrganisationID, created.OrganisationID)

	// The same email can't be registered twice.
	statusCode, _, err = r.MakeRequest(ctx, http.MethodPost, "v1/user", nil, payload, adminHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))
}

func TestUserCreateMasterAdminCanAssignRole(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	_, masterHeader := createUser(t, "master", userrole.MasterAdminID)

	payload := usermodel.User{
		Name:         "neworgadmin",
		EmailAddress: uniqueEmail("promoted"),
		Password:     testPassword,
		RoleID:       userrole.OrganisationAdminID,
	}

	statusCode, res, err := r.MakeRequest(ctx, http.MethodPost, "v1/user", nil, payload, masterHeader)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))
	assert.Equal(t, userrole.OrganisationAdminID, decodeUser(t, res).RoleID)
}

func TestDeactivatedUserLosesAccess(t *testing.T) {
	r, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.Background()

	u, header := createUser(t, "deactivated", "")

	statusCode, _, err := r.MakeRequest(ctx, http.MethodGet, "v1/me", nil, nil, header)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	app, err := testtools.CreateApp()
	assert.NoError(t, err)

	_, err = app.DB.Exec("UPDATE user SET is_active = 0 WHERE id = ?", u.ID)
	assert.NoError(t, err)

	// The existing session stops working straight away...
	statusCode, _, err = r.MakeRequest(ctx, http.MethodGet, "v1/me", nil, nil, header)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusUnauthorized))

	// ...and a new login is refused.
	login := usermodel.User{EmailAddress: u.EmailAddress, Password: testPassword}
	statusCode, _, err = r.MakeRequest(ctx, http.MethodPost, "v1/auth/login", nil, login)
	assert.NoError(t, err)
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusBadRequest))
}

func TestEmailAddressAlreadyExists(t *testing.T) {
	app, err := testtools.CreateApp()
	assert.NoError(t, err)

	srv := userservice.New(app)

	// A missing row is "doesn't exist", not an error.
	exists, err := srv.EmailAddressAlreadyExists(context.Background(), uniqueEmail("nobody"))
	assert.NoError(t, err)
	assert.False(t, exists)

	u, _ := createUser(t, "exists", "")

	exists, err = srv.EmailAddressAlreadyExists(context.Background(), u.EmailAddress)
	assert.NoError(t, err)
	assert.True(t, exists)
}
