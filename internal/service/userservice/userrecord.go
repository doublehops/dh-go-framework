package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/model/userrole"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// DefaultOrganisationID is the organisation seeded by migrations/20240923153040_create_user_table.up.sql.
// Users created without an organisation (e.g. test fixtures) are placed in it.
const DefaultOrganisationID = "00000000-0000-0000-0000-000000000100"

// UserRecord wraps a user model with service-level methods so callers don't
// need to pass the model back into the service for subsequent operations.
type UserRecord struct {
	Record *user.User
	svc    *UserService
}

// IsActive returns true if the user account is active.
func (u *UserRecord) IsActive() bool {
	return u.Record.IsActive == 1
}

// IsMasterAdmin returns true if the user holds the master_admin role.
func (u *UserRecord) IsMasterAdmin() bool {
	return u.Record.RoleID == userrole.MasterAdminID
}

// IsOrganisationAdmin returns true if the user holds the organisation_admin role.
func (u *UserRecord) IsOrganisationAdmin() bool {
	return u.Record.RoleID == userrole.OrganisationAdminID
}

// HasAdminRights returns true if the user holds any administrative role.
func (u *UserRecord) HasAdminRights() bool {
	return u.IsMasterAdmin() || u.IsOrganisationAdmin()
}

// HasAnyRole returns true if the user's global role ID matches any of the given role IDs.
func (u *UserRecord) HasAnyRole(roleIDs ...string) bool {
	for _, id := range roleIDs {
		if u.Record.RoleID == id {
			return true
		}
	}

	return false
}

// CheckPasswordHash returns true if input matches the stored hashed password.
func (u *UserRecord) CheckPasswordHash(input string) bool {
	return u.svc.CheckPasswordHash(input, u.Record.Password)
}

// CreateSession creates a new authenticated session for this user.
func (u *UserRecord) CreateSession(ctx context.Context) (*usersession.UserSession, error) {
	return u.svc.CreateUserSession(ctx, u.Record)
}

// HasPermissionToAccessUser returns true if the requesting user should have access to the given user.
func (u *UserRecord) HasPermissionToAccessUser(record *user.User) bool {
	return (u.IsOrganisationAdmin() && u.Record.OrganisationID == record.OrganisationID) ||
		u.Record.ID == record.ID ||
		u.IsMasterAdmin()
}

// Create persists a new user record after hashing the password.
func (u *UserRecord) Create(ctx context.Context) (*UserRecord, error) {
	if u.Record.OrganisationID == "" {
		u.Record.OrganisationID = DefaultOrganisationID
	}

	u.Record.IsActive = 1

	if u.Record.RoleID == "" {
		u.Record.RoleID = userrole.StandardUserID
	}

	hashedPassword, err := u.svc.HashPassword(u.Record.Password)
	if err != nil {
		u.svc.Log.Error(ctx, "unable to hash password", logga.KVPs{"error": err.Error()})

		return nil, errors.New("unable to hash password")
	}

	u.Record.Password = hashedPassword

	if err := u.Record.SetCreated(ctx); err != nil {
		u.svc.Log.Error(ctx, "error in SetCreated", logga.KVPs{"error": err.Error()})
	}

	tx := u.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err = u.svc.userRepo.Create(ctx, tx, u.Record); err != nil {
		u.svc.Log.Error(ctx, service.UnableToSaveRecord+" "+err.Error(), nil)

		return u, req.ErrCouldNotSaveRecord
	}

	if err = tx.Commit(); err != nil {
		u.svc.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	if err = u.svc.userRepo.GetByID(ctx, u.svc.DB, u.Record.ID, u.Record); err != nil {
		u.svc.Log.Error(ctx, service.UnableToRetrieveRecord+" "+err.Error(), nil)

		return u, errors.New(service.UnableToRetrieveRecord)
	}

	return u, nil
}

// Update saves changes to an existing user record.
func (u *UserRecord) Update(ctx context.Context, payload *user.User) error {
	u.Record.Name = payload.Name
	u.Record.SetUpdated(ctx)

	tx := u.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err := u.svc.userRepo.Update(ctx, tx, u.Record); err != nil {
		u.svc.Log.Error(ctx, service.UnableToUpdateRecord+" "+err.Error(), nil)
	}

	if err := tx.Commit(); err != nil {
		u.svc.Log.Error(ctx, "unable to commit transaction"+err.Error(), nil)
	}

	if err := u.svc.userRepo.GetByID(ctx, u.svc.DB, u.Record.ID, u.Record); err != nil {
		return errors.New(service.UnableToRetrieveRecord)
	}

	return nil
}

// Delete soft-deletes the user record.
func (u *UserRecord) Delete(ctx context.Context) error {
	u.Record.SetDeleted(ctx)

	tx := u.svc.DB.MustBegin()
	defer tx.Rollback() //nolint:errcheck

	if err := u.svc.userRepo.Delete(ctx, tx, u.Record); err != nil {
		u.svc.Log.Error(ctx, "unable to delete record. "+err.Error(), nil)

		return fmt.Errorf("unable to delete record")
	}

	if err := tx.Commit(); err != nil {
		u.svc.Log.Error(ctx, service.UnableToCommitTransaction+" "+err.Error(), nil)
	}

	return nil
}
