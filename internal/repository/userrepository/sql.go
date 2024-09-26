package userrepository

var insertRecordSQL = `INSERT INTO user (
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_expire,
	is_active,
	created_at,
	updated_at,
	deleted_at
	  ) VALUES (
	:organisation_id,
	:name,
	:email_address,
	:email_verified,
	:password,
	:password_reset_string,
	:password_reset_expire,
	:is_active,
	:created_at,
	:updated_at,
	:deleted_at
	)
`

var updateRecordSQL = `UPDATE user SET
	organisation_id=:organisation_id,
	name=:name,
	email_address=:email_address,
	email_verified=:email_verified,
	password=:password,
	password_reset_string=:password_reset_string,
	password_reset_expire=:password_reset_expire,
	is_active=:is_active,
	created_at=:created_at,
	updated_at=:updated_at,
	deleted_at=:deleted_at
	WHERE id=?
`

var deleteRecordSQL = `UPDATE user SET
    updated_at=?,
    deleted_at=?
	WHERE id=?
`

var selectByIDQuery = `SELECT 
	id,
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_expire,
	is_active,
	created_at,
	updated_at,
	deleted_at
    FROM user
    WHERE id=?
    AND deleted_at IS NULL`

var selectByEmailAddressQuery = `SELECT 
	id,
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_expire,
	is_active,
	created_at,
	updated_at,
	deleted_at
    FROM user
    WHERE email_address=?`

var selectCollectionQuery = `SELECT 
	id,
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_expire,
	is_active,
	created_at,
	updated_at,
	deleted_at
    FROM user
`

var selectCollectionCountQuery = `SELECT 
    COUNT(*) count
    FROM user
`
