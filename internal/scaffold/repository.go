package scaffold

import (
	"context"
	"fmt"
	"strings"
)

const (
	repositoryTemplate    = "./internal/scaffold/templates/repository.tmpl"
	repositorySQLTemplate = "./internal/scaffold/templates/repositorysql.tmpl"
)

// createRepository will create the repository.
func (s *Scaffold) createRepository(ctx context.Context, m Model) error {
	m.ModelStructProperties = getStructProperties(m.Columns)
	m.InsertFields, m.UpdateFields, m.ScanFields = s.getQueryFields(m.Columns)
	path := fmt.Sprintf("%s/%s/%s", s.pwd, s.Config.Paths.Repository, m.RepositoryName)
	repositoryFilename := fmt.Sprintf("%s/%s.go", path, m.RepositoryFilename)
	sqlFilename := fmt.Sprintf("%s/sql.go", path)

	err := MkDir(path)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	// Write repository file.
	err = s.writeFile(repositoryTemplate, repositoryFilename, m)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	s.l.Info(ctx, "repository has been written: "+repositoryFilename, nil)

	// Write repository SQL file.
	s.setColumnSQLParams(&m)

	err = s.writeFile(repositorySQLTemplate, sqlFilename, m)
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)

		return err
	}

	s.l.Info(ctx, "repository SQL has been written: "+repositoryFilename, nil)

	return nil
}

// getQueryFields will build a string for the various SQL queries for the given table.
func (s *Scaffold) getQueryFields(cols []column) (string, string, string) {
	var insertColumns []string
	var selectColumns []string
	var updateColumns []string

	for _, f := range cols {
		insertCol := fmt.Sprintf("model.%s", f.CapitalisedAbbr)
		insertColumns = append(insertColumns, insertCol)

		selectCol := fmt.Sprintf("&record.%s", f.CapitalisedAbbr)
		selectColumns = append(selectColumns, selectCol)

		// Don't add `id` column to beginning of update statement.
		if f.Original == "id" {
			continue
		}
		updateCol := fmt.Sprintf("model.%s", f.CapitalisedAbbr)
		updateColumns = append(updateColumns, updateCol)
	}

	updateColumns = append(updateColumns, "model.ID")

	insertFields := strings.Join(insertColumns[1:], ", ")
	updateFields := strings.Join(updateColumns, ", ")
	scanFields := strings.Join(selectColumns, ", ")

	return insertFields, updateFields, scanFields
}

// setColumnSQLParams will build the parameter count for each SQL query.
func (s *Scaffold) setColumnSQLParams(m *Model) {
	var (
		insertCols = ""
		insertQs   = ""
		updateStmt = ""
		selectStmt = ""
	)

	for _, col := range m.Columns {
		selectStmt += fmt.Sprintf("\t%s,\n", col.Original)

		if col.Original == "id" { // Don't include ID field in queries.
			continue
		}

		insertCols += fmt.Sprintf("\t%s,\n", col.Original)
		insertQs += "?,\n"
		updateStmt += fmt.Sprintf("\t%s=?,\n", col.Original)
	}

	// Remove two last chars (comma and carriage return) of each string.
	insertCols = insertCols[:len(insertCols)-2]
	insertQs = insertQs[:len(insertQs)-2]
	updateStmt = updateStmt[:len(updateStmt)-2]
	selectStmt = selectStmt[:len(selectStmt)-2]

	m.SQLCreate = insertCols
	m.SQLCreateQs = insertQs
	m.SQLUpdate = updateStmt
	m.SQLSelect = selectStmt
}
