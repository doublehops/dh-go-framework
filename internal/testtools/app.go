package testtools

import (
	"fmt"
	"log"
	"os"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/service"
)

func CreateApp() (*service.App, error) {
	cfg, err := config.New("./config_test.json")
	if err != nil {
		log.Printf("error starting main. %s", err.Error())
		os.Exit(1)
	}

	logCfg := &config.Logging{
		Writer:       "stdout",
		LogLevel:     "DEBUG",
		OutputFormat: "JSON",
	}

	l, err := logga.New(logCfg)
	if err != nil {
		return nil, fmt.Errorf("CreateTestUser: failed to create logger: %w", err)
	}

	database, err := db.New(l, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("CreateTestUser: failed to connect to database: %w", err)
	}

	app := &service.App{
		DB:  database,
		Log: l,
	}

	return app, nil
}
