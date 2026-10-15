package scaffold

import (
	"context"
	"fmt"
	"strings"

	"github.com/doublehops/go-common/str"
)

const modelTemplate = "./internal/scaffold/templates/model.tmpl"

// createModel will create the model.
func (s *Scaffold) createModel(ctx context.Context, m Model) error {
	s.l.Info(ctx, ">>>>>>>> in createModel()", nil)
	m.ModelStructProperties = getStructProperties(m.Columns)
	m.ValidationRules = s.getValidationRules(m)
	m.HasTime = hasTimeProperty(m.Columns)
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Paths.Model, m.LowerCase)
	filename := fmt.Sprintf("%s/%s.go", path, m.LowerCase)

	err := MkDir(path)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	err = s.writeFile(modelTemplate, filename, m)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	s.l.Info(ctx, "model has been written: "+filename, nil)

	return nil
}

// Column names given special treatment by the generator.
const (
	colID        = "id"
	colCreatedBy = "created_by"
	colCreatedAt = "created_at"
	colDeletedAt = "deleted_at"
)

// baseModelColumns are provided by model.BaseModel, so they're not generated on the model itself.
var baseModelColumns = []string{colID, colCreatedBy, "updated_by", colCreatedAt, "updated_at", colDeletedAt}

// hasTimeProperty reports whether any generated model property is a *time.Time.
func hasTimeProperty(columns []column) bool {
	for _, col := range columns {
		if col.Type == typeDatetime && !str.SliceContains(col.Original, baseModelColumns) {
			return true
		}
	}

	return false
}

// getStructProperties will build the struct properties.
func getStructProperties(columns []column) string {
	var properties string

	for _, col := range columns {
		if str.SliceContains(col.Original, baseModelColumns) {
			continue
		}

		properties += fmt.Sprintf("%s %s `json:\"%s\" db:\"%s\"`\n", col.CapitalisedAbbr, col.Type, col.CamelCase, col.Original)
	}

	return properties
}

// getValidationRules will build the validation rules defined by the column names.
func (s *Scaffold) getValidationRules(m Model) string {
	var rules string

	for _, col := range m.Columns {
		if str.SliceContains(col.Original, baseModelColumns) {
			continue
		}

		rules += getRule(col, m)
	}

	return rules
}

// getRule will return a starting validation rule based on the column type, to be refined by hand.
// Columns ending in _id that reference another table are good candidates for
// validator.RecordExistsByID. Types with no sensible default (datetime, bool) get no rule.
func getRule(col column, m Model) string {
	noLint := "//nolint:govet"

	switch col.Type {
	case typeString:
		return fmt.Sprintf("{\"%s\", %s.%s, true, []validator.ValidationFuncs{validator.LengthInRange(1, 255, \"\")}}, %s\n", col.CamelCase, m.FirstInitial, col.CapitalisedAbbr, noLint)
	case typeInt:
		return fmt.Sprintf("{\"%s\", %s.%s, true, []validator.ValidationFuncs{validator.IsInt(\"\")}}, %s\n", col.CamelCase, m.FirstInitial, col.CapitalisedAbbr, noLint)
	case typeBool, typeDatetime:
		return ""
	}

	return ""
}

// getPropertyType will check which column type the property is and return a corresponding
// Go Type to use in the model's struct.
func getPropertyType(propType string) columnType {
	if strings.Contains(propType, "int") {
		return typeInt
	}
	if strings.Contains(propType, "char") {
		return typeString
	}
	if strings.Contains(propType, "datetime") {
		return typeDatetime
	}

	return typeString
}
