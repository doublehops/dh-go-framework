// Package usersessionrepository handles database operations for user session records.
package usersessionrepository

var insertRecordSQL = `INSERT INTO user_session (
	user_id,
	token,
	last_request,
	created_at,
	updated_at,
	deleted_at
	) VALUES (
	:user_id,
	:token,
	:last_request,
	:created_at,
	:updated_at,
	:deleted_at
	)
`

var updateLastRequestSQL = `UPDATE user_session SET
	last_request=:last_request,
	updated_at=:updated_at
	WHERE id=:id
`

var deleteRecordSQL = `UPDATE user_session SET
	updated_at=:updated_at,
	deleted_at=:deleted_at
	WHERE id=:id
`

var selectByIDQuery = `SELECT
	id, user_id, token, last_request, created_at, updated_at, deleted_at
	FROM user_session
	WHERE id=?
	AND deleted_at IS NULL`

var selectByUserIDQuery = `SELECT
	id, user_id, token, last_request, created_at, updated_at, deleted_at
	FROM user_session
	WHERE user_id=?
	AND deleted_at IS NULL`

//nolint:gosec
var selectByTokenQuery = `SELECT
	id, user_id, token, last_request, created_at, updated_at, deleted_at
	FROM user_session
	WHERE token=?
	AND deleted_at IS NULL`

var selectCollectionQuery = `SELECT
	id, user_id, token, last_request, created_at, updated_at, deleted_at
	FROM user_session
`

var selectCollectionCountQuery = `SELECT
	COUNT(*) count
	FROM user_session
`
