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
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Config.Paths.Model, m.LowerCase)
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

// getStructProperties will build the struct properties.
func getStructProperties(columns []column) string {
	ignoreColumns := []string{"id", "created_at", "updated_at", "deleted_at"}

	var properties string

	for _, col := range columns {
		if str.SliceContains(col.Original, ignoreColumns) {
			continue
		}

		properties += fmt.Sprintf("%s %s `json:\"%s\" db:\"%s\"`\n", col.CapitalisedAbbr, col.Type, col.CamelCase, col.Original)
	}

	return properties
}

// getValidationRules will build the validation rules defined by the column names.
func (s *Scaffold) getValidationRules(m Model) string {
	var rules string

	ignoreColumns := []string{"id", "created_at", "updated_at", "deleted_at"}

	for _, col := range m.Columns {
		if str.SliceContains(col.Original, ignoreColumns) {
			continue
		}

		rules += getRule(col, m)
	}

	return rules
}

// getRule will return a validation rule based on the column type.
func getRule(col column, m Model) string {
	var rule string

	noLint := "//nolint:govet"

	switch col.Type {
	case "string":
		rule = fmt.Sprintf("{\"%s\", %s.%s, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, \"\")}}, %s\n", col.CamelCase, m.FirstInitial, col.CapitalisedAbbr, noLint)
	case "int":
		rule = fmt.Sprintf("{\"%s\", %s.%s, true, []validator.ValidationFuncs{validator.IsInt(\"\")}}, %s\n", col.CamelCase, m.FirstInitial, col.CapitalisedAbbr, noLint)
	default:
		rule = fmt.Sprintf("{\"%s\", %s.%s, true, []validator.ValidationFuncs{validator.LengthInRange(3, 8, \"\")}}, %s,\n", col.CamelCase, m.FirstInitial, col.CapitalisedAbbr, noLint)
	}

	return rule
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
