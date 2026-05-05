package testtools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	migrate "github.com/doublehops/dh-go-framework/internal/migration"
)

// RunMigrations runs all pending "up" migrations against the database configured in cfg.
// It locates the migrations directory by walking up from the current working directory
// until it finds the project root (identified by go.mod).
func RunMigrations(cfg *config.Config) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}

	logCfg := &config.Logging{
		Writer:       "stdout",
		LogLevel:     "DEBUG",
		OutputFormat: "JSON",
	}

	l, err := logga.New(logCfg)
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to create logger: %w", err)
	}

	database, err := db.New(l, cfg.DB)
	if err != nil {
		return fmt.Errorf("RunMigrations: failed to connect to database: %w", err)
	}

	action := &migrate.Action{
		Action: "up",
		Number: 9999,
		DB:     database,
		Path:   filepath.Join(projectRoot, "migrations"),
	}

	if err := action.Migrate(); err != nil {
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
