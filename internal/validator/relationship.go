package validator

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
)

// RecordExistsDefaultMessage is the default error message when a referenced record can't be found.
const RecordExistsDefaultMessage = "referenced record does not exist"

// dbConn is the connection RecordExistsByID checks against. ValidationFuncs' signature
// (func(bool, interface{}) (bool, string)) has no room to thread a DB handle through per-call, so
// it's wired up once at startup via Init instead — see cmd/server/run.go.
var dbConn *sqlx.DB

// Init wires up the database connection used by DB-backed validators such as RecordExistsByID.
// Must be called once during application startup, before any request is validated.
func Init(conn *sqlx.DB) {
	dbConn = conn
}

// RecordExistsByID returns a ValidationFuncs that checks a (non-soft-deleted) record with the given
// ID exists in the table matching model's type name, e.g. author.Author{} looks in "author".
func RecordExistsByID(model any, errorMessage string) ValidationFuncs {
	table := tableNameFromModel(model)

	return func(required bool, value interface{}) (bool, string) {
		if errorMessage == "" {
			errorMessage = RecordExistsDefaultMessage
		}

		var v string
		var ok bool

		if v, ok = value.(string); !ok {
			return false, ProcessingPropertyError
		}

		if v == "" && required {
			return false, RequiredPropertyError
		}

		if v == "" && !required {
			return true, ""
		}

		if dbConn == nil {
			return false, ProcessingPropertyError
		}

		var count int

		query := "SELECT COUNT(*) FROM " + table + " WHERE id = ? AND deleted_at IS NULL"
		if err := dbConn.Get(&count, query, v); err != nil {
			return false, ProcessingPropertyError
		}

		if count == 0 {
			return false, errorMessage
		}

		return true, ""
	}
}

// tableNameFromModel derives a table name from a model's Go type name — e.g. author.Author{}
// becomes "author", widgetpart.WidgetPart{} becomes "widget_part".
// This mirrors every model in internal/model today: table names are the snake_case of the struct name.
func tableNameFromModel(model any) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	var snake strings.Builder

	for i, r := range t.Name() {
		if unicode.IsUpper(r) && i > 0 {
			snake.WriteRune('_')
		}

		snake.WriteRune(unicode.ToLower(r))
	}

	return snake.String()
}
