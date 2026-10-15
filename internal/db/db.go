// Package db provides database connection setup and management.
package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/logga"
)

const (
	maxAttempts = 3
	delay       = 5
)

// New creates and returns a new database connection, retrying up to maxAttempts times.
func New(l *logga.Logga, cfg config.DB) (*sqlx.DB, error) {
	l.Log.Info("opening database connection")

	var db *sqlx.DB
	var err error

	// FormatDSN escapes the credentials, so a password containing '@', '/' or ':' still connects.
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.User
	mysqlCfg.Passwd = cfg.Pass
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%s", cfg.Host, cfg.PortOrDefault())
	mysqlCfg.DBName = cfg.Name
	mysqlCfg.ParseTime = true
	connectStr := mysqlCfg.FormatDSN()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err = sqlx.Connect("mysql", connectStr)
		if err == nil {
			return db, nil
		}

		l.Log.Error(fmt.Sprintf("unable to open db connection. %s", err))
		// cfg.Pass is deliberately not logged: logs are commonly shipped to external storage, and a
		// wrong host or database name is what actually needs diagnosing here.
		l.Log.Error(fmt.Sprintf("Database connection error. user: %s; host: %s; name: %s", cfg.User, cfg.Host, cfg.Name))

		time.Sleep(delay * time.Second)
	}

	return nil, fmt.Errorf("unable to open database connection. max attempts reached. Attempts: %d", maxAttempts)
}
