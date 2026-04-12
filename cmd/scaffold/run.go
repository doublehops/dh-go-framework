// Package main is the entry point for the scaffold code generation tool.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/scaffold"
)

/*
  This file just serves as an example of how you would add this library to your project.
*/

func main() {
	if err := run(); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}

// run will begin the scaffolding.
func run() error {
	ctx := context.Background()

	var tableName string

	flag.StringVar(&tableName, "table", "", "the database table (and model) to scaffold")

	configFile := flag.String("config", "config.json", "Scaffold file to use")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("there was an error with os.Getwd(). %s", err.Error())
	}

	// Setup config.
	cfg, err := config.New(*configFile)
	if err != nil {
		return fmt.Errorf("error starting. %s", err.Error())
	}

	// Setup logger.
	l, err := logga.New(&cfg.Logging)
	if err != nil {
		return fmt.Errorf("error configuring logger. %s", err.Error())
	}

	l.Info(context.TODO(), "flags", logga.KVPs{"table": tableName})

	// Setup db connection.
	DB, err := db.New(l, cfg.DB)
	if err != nil {
		errMsg := fmt.Sprintf("error creating database connection. %s", err.Error())
		l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	scf, err := GetScaffoldConfig(pwd)
	if err != nil {
		errMsg := fmt.Sprintf("error getting scaffold config. %s", err.Error())
		l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	s := scaffold.New(pwd, scf, tableName, DB, l)
	err = s.Run()
	if err != nil {
		return err
	}

	return nil
}

// GetScaffoldConfig will pull the config from the config file.
func GetScaffoldConfig(pwd string) (scaffold.Config, error) {
	cp := scaffold.Config{}

	relPath := pwd + "/cmd/scaffold/config.json"

	f, err := os.ReadFile(relPath) //nolint:gosec
	if err != nil {
		log.Printf("unable to read config file - %s. %s", relPath, err.Error())

		return cp, fmt.Errorf("unable to read config file `%s`. %w", relPath, err)
	}

	if err := json.Unmarshal(f, &cp); err != nil {
		return cp, err
	}

	return cp, nil
}
