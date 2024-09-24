package userrepository

var insertRecordSQL = `INSERT INTO user (
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_timeout,
	is_active,
	created_at,
	updated_at,
	deleted_at
	  ) VALUES (
?,
?,
?,
?,
?,
?,
?,
?,
?,
?,
?
	)
`

var updateRecordSQL = `UPDATE user SET
	organisation_id=?,
	name=?,
	email_address=?,
	email_verified=?,
	password=?,
	password_reset_string=?,
	password_reset_timeout=?,
	is_active=?,
	created_at=?,
	updated_at=?,
	deleted_at=?
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
	password_reset_timeout,
	is_active,
	created_at,
	updated_at,
	deleted_at
    FROM user
    WHERE id=?
    AND deleted_at IS NULL`

var selectCollectionQuery = `SELECT 
	id,
	organisation_id,
	name,
	email_address,
	email_verified,
	password,
	password_reset_string,
	password_reset_timeout,
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
