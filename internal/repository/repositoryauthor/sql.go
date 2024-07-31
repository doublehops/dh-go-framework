package repositoryauthor

var insertRecordSQL = `INSERT INTO author (
	user_id,
	name,
    created_by,
    updated_by,
    created_at,
    updated_at
	  ) VALUES (
	:user_id,
	:name,
	:created_by,
	:updated_by,
	:created_at,
	:updated_at
	)
`

var updateRecordSQL = `UPDATE author SET 
	name=?,
    updated_by=?,
    updated_at=?
	WHERE id=?
`

var deleteRecordSQL = `UPDATE author SET 
    updated_by=?,
    updated_at=?,
    deleted_at=?
	WHERE id=?
`

var selectByIDQuery = `SELECT 
    *
    FROM author
    WHERE id=?
    AND deleted_at IS NULL`

var selectCollectionQuery = `SELECT 
    id,
    user_id,
    name,
    created_by,
    updated_by,
    created_at,
    updated_at
    FROM author
`

var selectCollectionCountQuery = `SELECT 
    COUNT(*) count
    FROM author
`
