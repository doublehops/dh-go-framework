// Package main is the entry point for the API server.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/routes"
	"github.com/doublehops/dh-go-framework/internal/runflags"
	"github.com/doublehops/dh-go-framework/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	flags := runflags.GetFlags()

	// Setup config.
	cfg, err := config.New(flags.ConfigFile)
	if err != nil {
		return fmt.Errorf("error starting main. %s", err.Error())
	}

	// Setup logger.
	l, err := logga.New(&cfg.Logging)
	if err != nil {
		return fmt.Errorf("error configuring logger. %s", err.Error())
	}

	lMsg := fmt.Sprintf("MySQL config found. %s", cfg.DB.Host)
	l.Info(ctx, lMsg, nil)

	// Setup db connection.
	DB, err := db.New(l, cfg.DB)
	if err != nil {
		return fmt.Errorf("error creating database connection. %s", err.Error())
	}

	App := &service.App{
		DB:  DB,
		Log: l,
	}

	router := httprouter.New()
	rts := routes.GetV1Routes(App)

	l.Info(ctx, "Adding routes", nil)
	for _, r := range rts.Routes() {
		log.Printf(">>> %s %s\n", r.Method(), r.Path())
		router.Handle(r.Method(), r.Path(), r.Handler())
	}

	mux := http.TimeoutHandler(router, time.Second*1, "Timeout!")

	l.Info(ctx, "Starting server on port :"+cfg.Host.Port, nil)

	srv := &http.Server{
		Addr:         ":" + cfg.Host.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err = srv.ListenAndServe(); err != nil {
		return fmt.Errorf("unable to start server. %s", err.Error())
	}

	return nil
}
