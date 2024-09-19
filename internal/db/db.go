package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Importing the MySQL driver for database/sql.
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/logga"
)

func New(l *logga.Logga, cfg config.DB) (*sqlx.DB, error) {
	l.Log.Info("opening database connection")

	connectStr := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", cfg.User, cfg.Pass, cfg.Host, cfg.Name)
	db, err := sqlx.Connect("mysql", connectStr)
	if err != nil {
		l.Log.Error(fmt.Sprintf("unable to open db connection. %s", err))
		l.Log.Error(fmt.Sprintf("Database connection error. user: %s; password: %s; host: %s; name: %s", cfg.User, cfg.Pass, cfg.Host, cfg.Name))

		return db, err
	}

	l.Log.Info(">>>>>> Database connection success")

	return db, nil
}
