// Package main is the entry point for the database migration CLI tool.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gm "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/doublehops/dh-go-framework/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}

func run() error {
	action := flag.String("action", "", "action: up, down, create")
	name := flag.String("name", "", "migration name (required for create)")
	number := flag.Int("number", 0, "number of migrations to run (0 = all for up, 1 for down)")
	configFile := flag.String("config", "config.json", "config file")
	flag.Parse()

	if *action == "" {
		fmt.Fprintln(os.Stderr, "Usage: -action [up|down|create] [-number N] [-name <name>]")
		os.Exit(1)
	}

	if *action == "create" {
		if *name == "" {
			fmt.Fprintln(os.Stderr, "-name is required for the create action")
			os.Exit(1)
		}

		return createMigration(*name)
	}

	cfg, err := config.New(*configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	m, cleanup, err := newMigrator(cfg)
	if err != nil {
		return err
	}

	defer cleanup()

	return runAction(m, *action, *number)
}

func runAction(m *gm.Migrate, action string, number int) error {
	version, dirty, _ := m.Version() //nolint:errcheck
	log.Printf("current version: %d (dirty: %v)", version, dirty)

	var err error

	switch action {
	case "up":
		if number == 0 {
			err = m.Up()
		} else {
			err = m.Steps(number)
		}
	case "down":
		n := number
		if n == 0 {
			n = 1
		}

		err = m.Steps(-n)
	default:
		return fmt.Errorf("unknown action %q — use up, down, or create", action)
	}

	if err == gm.ErrNoChange {
		log.Println("no migrations to run")

		return nil
	}

	return err
}

func newMigrator(cfg *config.Config) (*gm.Migrate, func(), error) {
	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true&multiStatements=true",
		cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Name)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening database: %w", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("error creating mysql driver: %w", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("error getting working directory: %w", err)
	}

	m, err := gm.NewWithDatabaseInstance("file://"+dir+"/migrations", "mysql", driver)
	if err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("error creating migrator: %w", err)
	}

	return m, func() { _ = sqlDB.Close() }, nil
}

func createMigration(name string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}

	ts := time.Now().Format("20060102150405")
	base := fmt.Sprintf("%s/migrations/%s_%s", dir, ts, name)

	if err := os.WriteFile(base+".up.sql", []byte("-- Write your up migration SQL here\n"), 0o600); err != nil {
		return fmt.Errorf("error creating up migration file: %w", err)
	}
	if err := os.WriteFile(base+".down.sql", []byte("-- Write your rollback SQL here\n"), 0o600); err != nil {
		return fmt.Errorf("error creating down migration file: %w", err)
	}

	log.Printf("created:\n  %s.up.sql\n  %s.down.sql", base, base)

	return nil
}
