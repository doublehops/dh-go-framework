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

// RunMigrations drops all tables and re-applies every migration so each test
// run starts from a known-clean state.
func RunMigrations(cfg *config.Config) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}

	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true&multiStatements=true",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	migrationsPath := "file://" + filepath.Join(projectRoot, "migrations")

	// openMigrator opens its own *sql.DB so that m.Close() (which closes the DB)
	// does not affect a subsequent migrator instance.
	openMigrator := func() (*gm.Migrate, error) {
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to open db: %w", err)
		}
		drv, err := mysql.WithInstance(sqlDB, &mysql.Config{})
		if err != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to create driver: %w", err)
		}
		m, err := gm.NewWithDatabaseInstance(migrationsPath, "mysql", drv)
		if err != nil {
			sqlDB.Close()
			return nil, fmt.Errorf("failed to create migrator: %w", err)
		}
		return m, nil
	}

	// Drop all tables so the test run starts from scratch.
	m, err := openMigrator()
	if err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}
	if err := m.Drop(); err != nil {
		_, _ = m.Close()
		return fmt.Errorf("RunMigrations: drop failed: %w", err)
	}
	_, _ = m.Close()

	// Fresh connection so the driver re-initialises schema_migrations, then apply all migrations.
	m, err = openMigrator()
	if err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}
	defer func() { _, _ = m.Close() }()

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
