package gomigration

// QuerySeparator is the delimiter used to split multiple SQL statements within a migration file.
const QuerySeparator = "------------------"

// GetLatestMigrationSQL and related vars are the SQL statements used to manage the migrations table.
var (
	GetLatestMigrationSQL = `SELECT * FROM migrations
	ORDER BY id DESC
	LIMIT 1
`

	CreateMigrationsTable = `CREATE TABLE migrations (
	id INT(11) NOT NULL AUTO_INCREMENT,
	filename VARCHAR(255),
	created_at DATETIME,
	PRIMARY KEY(id)
)`

	CheckMigrationsTableExistsSQL = `SHOW TABLES`

	InsertMigrationRecordIntoTableSQL = `INSERT INTO migrations
	(filename,created_at)
	VALUES
	(?,NOW())
`

	RemoveMigrationRecordFromTableSQL = `DELETE FROM migrations
	WHERE filename = ?
`
)
