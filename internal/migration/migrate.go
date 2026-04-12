package gomigration

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql" // nolint:revive
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/migration/helpers"
)

// Action holds the configuration for a migration run including direction, number, and DB connection.
type Action struct {
	Action string
	Number int
	Name   string

	DB   *sqlx.DB
	Path string
}

// File represents a migration file with its parsed SQL queries.
type File struct {
	Filename string
	Queries  []string
}

// TrimExtension removes the .up.sql or .down.sql suffix from a migration filename.
func TrimExtension(filename string) string {
	var str string

	str = strings.Replace(filename, ".up.sql", "", 1)
	str = strings.Replace(str, ".down.sql", "", 1)

	return str
}

// TableList is a slice of Table records used when inspecting database tables.
type TableList []Table

// Table represents a database table name, used to check for the migrations table.
type Table struct {
	Name string
}

// Migrate runs the migration action (create, up, or down) as specified in a.Action.
//
//nolint:cyclop
func (a *Action) Migrate() error {
	var err error

	err = a.ensureMigrationsTableExists()
	if err != nil {
		return err
	}

	if a.Action == "create" {
		err = a.CreateMigration(a.Path)

		return err
	}

	if a.Action == "up" {
		pendingFiles, err := a.getPendingMigrationFiles()
		if err != nil {
			return err
		}
		if len(pendingFiles) == 0 {
			helpers.PrintMsg("There are no pending migrations\n")

			return nil
		}
		migrationFiles, err := a.parseMigrations(pendingFiles)
		if err != nil {
			return err
		}

		if err := a.MigrateUp(migrationFiles); err != nil {
			return err
		}
	}

	if a.Action == "down" {
		previousFiles, err := a.getMigrationFilesToRollBack()
		if err != nil {
			return err
		}
		if len(previousFiles) == 0 {
			helpers.PrintMsg("There are no previous migrations to rollback\n")

			return nil
		}
		migrationFiles, err := a.parseMigrations(previousFiles)
		if err != nil {
			return err
		}

		if err := a.MigrateDown(migrationFiles); err != nil {
			return err
		}
	}

	return nil
}

// IsValidAction returns true if the key is one of the supported migration actions.
func (a *Action) IsValidAction(key string) bool {
	validActions := []string{
		"create",
		"up",
		"down",
	}

	for _, item := range validActions {
		if item == key {
			return true
		}
	}

	return false
}

// PrintHelp prints usage instructions for the migration CLI tool and exits.
func (a *Action) PrintHelp() {
	helpMsg := `
Usage: <your_script> -action=<action> -number=<number>
Examples: 
./main.go -action create -name add_user_table // Will create a new migration file with template.
./main.go -action up -number 1 // number is optional. Will run all migrations if not included.
./main.go -action down -number 1 // number is optional. Will run only one migration if not included.
`
	_, _ = os.Stderr.WriteString(helpMsg)
	os.Exit(1)
}

// ensureMigrationsTableExists to create table to track migrations.
func (a *Action) ensureMigrationsTableExists() error {
	var tableList TableList
	rows, err := a.DB.Query(CheckMigrationsTableExistsSQL)
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()
	if err != nil {
		return fmt.Errorf("++++>>>> Error: %w", err)
	}

	for rows.Next() {
		var t Table
		if err := rows.Scan(&t.Name); err != nil {
			return err
		}
		tableList = append(tableList, t)
	}

	for _, tbl := range tableList {
		if tbl.Name == "migrations" { // Already exists
			return nil
		}
	}

	err = a.addMigrationTable()
	if err != nil {
		return err
	}

	return nil
}
