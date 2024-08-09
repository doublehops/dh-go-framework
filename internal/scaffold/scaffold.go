package scaffold

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
)

const goModuleFile = "./go.mod"

type Scaffold struct {
	DB  *sqlx.DB
	l   *logga.Logga
	pwd string

	tableName string
	Config
}

type Config struct {
	Paths Paths `json:"paths"`
}

type Paths struct {
	Handlers   string `json:"handlers"`
	Model      string `json:"model"`
	Repository string `json:"repository"`
	Service    string `json:"service"`
}

type Model struct {
	Name           string
	FirstInitial   string
	CamelCase      string
	PascalCase     string
	SnakeCase      string
	LowerCase      string
	KebabCase      string
	Initialisation string
	Module         string

	ServiceFilename    string
	RepositoryFilename string

	ServiceName    string
	RepositoryName string

	ModelStructProperties string
	ValidationRules       string
	InsertFields          string
	UpdateFields          string
	ScanFields            string

	Columns []column

	SQLCreate   string
	SQLCreateQs string
	SQLUpdate   string
	SQLSelect   string
}

type column struct {
	Original        string
	Type            columnType
	LowerCase       string
	PascalCase      string
	CamelCase       string
	CapitalisedAbbr string
}

type columnType string

const (
	typeInt      columnType = "int"
	typeString   columnType = "string"
	typeBool     columnType = "bool"
	typeDatetime columnType = "*datetime"
)

func New(pwd string, cfg Config, tableName string, db *sqlx.DB, logga *logga.Logga) *Scaffold {
	return &Scaffold{
		pwd:       pwd,
		DB:        db,
		l:         logga,
		tableName: tableName,
		Config:    cfg,
	}
}

func (s *Scaffold) Run() error {
	ctx := context.Background()

	columns, err := s.getTableDefinition()
	if err != nil {
		s.l.Error(ctx, "error getting column. "+err.Error(), nil)
		return errors.New("failed to run. " + err.Error())
	}

	moduleName, err := getModuleName()
	if err != nil {
		s.l.Error(ctx, err.Error(), nil)
		return errors.New("failed to run. " + err.Error())
	}

	ms := Model{
		Name:           s.tableName,
		FirstInitial:   GetFirstRune(s.tableName),
		CamelCase:      ToCamelCase(s.tableName),
		PascalCase:     ToPascalCase(s.tableName),
		KebabCase:      ToKebabCase(s.tableName),
		SnakeCase:      s.tableName,
		LowerCase:      RemoveUnderscores(s.tableName),
		Initialisation: ToInitialisation(s.tableName),
		Module:         moduleName,

		ServiceFilename:    "service" + RemoveUnderscores(s.tableName) + ".go",
		RepositoryFilename: "repository" + RemoveUnderscores(s.tableName) + ".go",

		ServiceName:    ToPascalCase(s.tableName) + "Service",
		RepositoryName: ToPascalCase(s.tableName) + "Repository",

		Columns: getColumnDefinitions(columns),
	}

	// Create model.
	err = s.createModel(ctx, ms)
	if err != nil {
		return err
	}

	// Create repository.
	err = s.createRepository(ctx, ms)
	if err != nil {
		return err
	}

	// Create service.
	err = s.createService(ctx, ms)
	if err != nil {
		return err
	}

	// Create handler.
	err = s.createHandler(ctx, ms)
	if err != nil {
		return err
	}

	// Create handler.
	err = s.printRoutes(ctx, ms)
	if err != nil {
		return err
	}

	return nil
}

// getModuleName will get the module name from go.mod to use to populate the templates.
func getModuleName() (string, error) {
	f, err := os.Open(goModuleFile)
	if err != nil {
		return "", errors.New("Opening go.mod failed. " + err.Error())
	}
	rawBytes, err := io.ReadAll(f)
	if err != nil {
		return "", errors.New("could not open module file" + err.Error())
	}
	lines := strings.Split(string(rawBytes), "\n")

	module := strings.Replace(lines[0], "module ", "", 1)

	return module, nil
}

type ColumnDefinition struct {
	columnName string
	columnType string
}

// getTableDefinition will get the table and column definitions which will be used to build all files.
func (s *Scaffold) getTableDefinition() ([]ColumnDefinition, error) {
	cols := []ColumnDefinition{}

	q := "DESCRIBE " + s.tableName
	rows, err := s.DB.Query(q)
	if err != nil || rows.Err() != nil {
		return cols, errors.New("error executing query. " + err.Error())
	}
	if rows.Err() != nil {
		return cols, errors.New("error in rows.Err(). " + err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.New("Error getting columns. " + err.Error())
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		columnName := fmt.Sprintf("%s", values[0])
		columnType := fmt.Sprintf("%s", values[1])
		cols = append(cols, ColumnDefinition{columnName: columnName, columnType: columnType})
	}

	return cols, nil
}

// getColumnDefinitions will get the definitions of each column.

func getColumnDefinitions(columns []ColumnDefinition) []column {
	var cols []column

	for _, col := range columns {
		col := column{
			Original:        col.columnName,
			Type:            getPropertyType(col.columnType),
			LowerCase:       strings.ToLower(ToCamelCase(col.columnName)),
			PascalCase:      ToPascalCase(col.columnName),
			CamelCase:       ToCamelCase(col.columnName),
			CapitalisedAbbr: CapitaliseAbbr(ToPascalCase(col.columnName)),
		}

		cols = append(cols, col)
	}

	return cols
}
