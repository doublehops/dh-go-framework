package testtools

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	gm "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/doublehops/dh-go-framework/internal/config"
)

// RunMigrations runs all pending "up" migrations against the database configured in cfg.
// It locates the migrations directory by walking up from the current working directory
// until it finds the project root (identified by go.mod).
func RunMigrations(cfg *config.Config) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true&multiStatements=true",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to open db: %w", err)
	}
	defer sqlDB.Close()

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to create driver: %w", err)
	}

	m, err := gm.NewWithDatabaseInstance(
		"file://"+filepath.Join(projectRoot, "migrations"),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != gm.ErrNoChange {
		return fmt.Errorf("RunMigrations: migration failed: %w", err)
	}

	return nil
}

// findProjectRoot walks up from the current working directory until it finds a go.mod file.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get working directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no go.mod found)")
		}

		dir = parent
	}
}
