// Package userrepository provides database operations for user records.
package userrepository

var insertRecordSQL = `INSERT INTO user (
	id,
	organisation_id,
	role_id,
	name,
	email_address,
	password,
	is_active,
	created_at,
	updated_at,
	deleted_at
	  ) VALUES (
	:id,
	:organisation_id,
	:role_id,
	:name,
	:email_address,
	:password,
	:is_active,
	:created_at,
	:updated_at,
	:deleted_at
	)
`

var updateRecordSQL = `UPDATE user SET
	name=:name,
	updated_by=:updated_by,
	updated_at=:updated_at
	WHERE id=:id
`

var deleteRecordSQL = `UPDATE user SET
	updated_by=:updated_by,
	updated_at=:updated_at,
	deleted_at=:deleted_at
	WHERE id=:id
`

var selectByIDQuery = `SELECT 
	*
    FROM user
    WHERE id=?
    AND deleted_at IS NULL`

var selectByEmailAddressQuery = `SELECT 
	*
    FROM user
    WHERE email_address=?
    AND deleted_at IS NULL`

var selectCollectionQuery = `SELECT 
	*
    FROM user
`

var selectCollectionCountQuery = `SELECT 
    COUNT(*) count
    FROM user
`
