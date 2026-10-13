// Package userrole defines the UserRole model for system-level role definitions.
package userrole

import "github.com/doublehops/dh-go-framework/internal/model"

const (
	// MasterAdminID is the hard-coded ID for the master admin role.
	MasterAdminID = "4ad2b30f-81d7-4b58-b950-db5c51e9590e"
	// OrganisationAdminID is the hard-coded ID for the organisation admin role.
	OrganisationAdminID = "6a8c21d8-201d-4420-a964-54a2b99c3db9"
	// StandardUserID is the hard-coded ID for the baseline "no special global role" role,
	// and the default user.role_id (see migrations/20240923153040_create_user_table.up.sql).
	StandardUserID = "6649c54a-4641-4554-9ec3-497192ac3df3"
)

// UserRole represents a global (system-level) role that can be assigned to a user.
type UserRole struct {
	model.BaseModel
	Name string `json:"name" db:"name"`
}
