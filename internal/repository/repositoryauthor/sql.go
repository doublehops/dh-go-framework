package repositoryauthor

var insertRecordSQL = `INSERT INTO {{.Name}} (
	user_id,
	name,
    created_by,
    updated_by,
    created_at,
    updated_at
	  ) VALUES (
	?,
	?,
	?,
	?,
	?,
	?
	)
`

var updateRecordSQL = `UPDATE {{.Name}} SET 
	name=?,
    updated_by=?,
    updated_at=?
	WHERE id=?
`

var deleteRecordSQL = `UPDATE {{.Name}} SET 
    updated_by=?,
    updated_at=?,
    deleted_at=?
	WHERE id=?
`

var selectByIDQuery = `SELECT 
    id,
    user_id,
    name,
    created_by,
    updated_by,
    created_at,
    updated_at
    FROM {{.Name}}
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
    FROM {{.Name}}
`

var selectCollectionCountQuery = `SELECT 
    COUNT(*) count
    FROM {{.Name}}
`
