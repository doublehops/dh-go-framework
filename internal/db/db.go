package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Importing the MySQL driver for database/sql.
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/logga"
)

const (
	maxAttempts = 5
	delay       = 5
)

func New(l *logga.Logga, cfg config.DB) (*sqlx.DB, error) {
	l.Log.Info("opening database connection")

	var db *sqlx.DB
	var err error

	connectStr := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sqlx.Connect("mysql", connectStr)
		if err == nil {
			return db, nil
		}

		l.Log.Error(fmt.Sprintf("unable to open db connection. %s", err))
		l.Log.Error(fmt.Sprintf("Database connection error. user: %s; password: %s; host: %s; name: %s", cfg.User, cfg.Pass, cfg.Host, cfg.Name))

		time.Sleep(delay * time.Second)
	}

	return nil, fmt.Errorf("unable to open database connection. max attempts reached. Attempts: %d", maxAttempts)
}
