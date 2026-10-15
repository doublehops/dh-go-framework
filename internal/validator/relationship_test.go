package validator

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
)

// Author is a local double for author.Author — RecordExistsByID only reflects on the
// type's name, so this avoids validator depending on any model package (which would in
// turn import validator, creating an import cycle).
type Author struct{}

// WidgetPart exercises the snake_case conversion of a multi-word type name.
type WidgetPart struct{}

// This test needs a migrated test database (see `make migrate_test`); the handler test packages
// also migrate it on startup.

func TestRecordExistsByID(t *testing.T) {
	cfg, err := config.New("../../config_test.json")
	assert.NoError(t, err)

	l, err := logga.New(&cfg.Logging)
	assert.NoError(t, err)

	conn, err := db.New(l, cfg.DB)
	assert.NoError(t, err)

	// Before Init the rule must fail safely rather than panic on a nil connection.
	dbConn = nil
	res, errMsg := RecordExistsByID(Author{}, "")(true, uuid.NewString())
	assert.False(t, res)
	assert.Equal(t, ProcessingPropertyError, errMsg)

	Init(conn)

	id := uuid.NewString()
	_, err = conn.Exec(
		"INSERT INTO author (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		id, "relationship-test-author",
	)
	assert.NoError(t, err)

	deletedID := uuid.NewString()
	_, err = conn.Exec(
		"INSERT INTO author (id, name, created_at, updated_at, deleted_at) VALUES (?, ?, NOW(), NOW(), NOW())",
		deletedID, "relationship-test-author-deleted",
	)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		required       bool
		value          interface{}
		expectedErrMsg string
		expectedResult bool
	}{
		{"existingRecord", true, id, "", true},
		{"softDeletedRecord", true, deletedID, RecordExistsDefaultMessage, false},
		{"nonExistentRecord", true, uuid.NewString(), RecordExistsDefaultMessage, false},
		{"emptyButNotRequired", false, "", "", true},
		{"emptyButRequired", true, "", RequiredPropertyError, false},
		{"notAString", true, 123, ProcessingPropertyError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := RecordExistsByID(Author{}, "")
			res, errMsg := f(tt.required, tt.value)
			assert.Equal(t, tt.expectedErrMsg, errMsg, "error message not as expected")
			assert.Equal(t, tt.expectedResult, res, "result not as expected")
		})
	}
}

func TestTableNameFromModel(t *testing.T) {
	tests := []struct {
		name     string
		model    any
		expected string
	}{
		{"value type", Author{}, "author"},
		{"pointer type", &Author{}, "author"},
		{"multi-word type", WidgetPart{}, "widget_part"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tableNameFromModel(tt.model))
		})
	}
}
