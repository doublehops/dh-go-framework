// Package main is the entry point for the API server.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/middleware"
	"github.com/doublehops/dh-go-framework/internal/routes"
	"github.com/doublehops/dh-go-framework/internal/runflags"
	"github.com/doublehops/dh-go-framework/internal/service"
	"github.com/doublehops/dh-go-framework/internal/validator"
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

	l.Info(ctx, fmt.Sprintf("MySQL config found. %s", cfg.DB.Host), nil)

	// Setup db connection.
	DB, err := db.New(l, cfg.DB)
	if err != nil {
		return fmt.Errorf("error creating database connection. %s", err.Error())
	}

	// DB-backed validation rules (e.g. validator.RecordExistsByID) check against this connection.
	validator.Init(DB)

	App := &service.App{
		DB:  DB,
		Log: l,
	}

	router := httprouter.New()
	rts := routes.GetV1Routes(App, cfg)

	l.Info(ctx, "Adding routes", nil)
	for _, r := range rts.Routes() {
		log.Printf(">>> %s %s\n", r.Method(), r.Path())
		router.Handle(r.Method(), r.Path(), r.Handler())
	}

	// requestTimeout bounds total handler execution (decode/process — not just routing).
	const requestTimeout = 30 * time.Second

	mux := http.TimeoutHandler(router, requestTimeout, "Timeout!")
	if cfg.Host.CORSOrigin != "" {
		mux = middleware.CORSMiddleware(cfg.Host.CORSOrigin, mux)
	}

	// A configured port of "0" asks the OS for any free port — used by the test harness
	// (internal/testtools/server.go) so parallel test packages never collide on a fixed port.
	// net.Listen (rather than srv.ListenAndServe, which listens internally) is required so the
	// actual assigned port can be recovered and reported below.
	listener, err := net.Listen("tcp", ":"+cfg.Host.Port)
	if err != nil {
		return fmt.Errorf("unable to start server. %s", err.Error())
	}

	actualPort := listener.Addr().(*net.TCPAddr).Port //nolint:forcetypeassert // net.Listen("tcp", ...) always yields a *TCPAddr

	l.Info(ctx, fmt.Sprintf("Starting server on port :%d", actualPort), nil)

	if cfg.Host.Port == "0" {
		// Plain, undecorated marker line written directly to stdout (not routed through the
		// logger, whose format varies) so internal/testtools/server.go can reliably parse the
		// OS-assigned port back out of the subprocess's captured stdout.
		_, _ = fmt.Fprintf(os.Stdout, "SERVER_LISTENING_PORT=%d\n", actualPort)
	}

	// ReadTimeout and WriteTimeout are deliberately unset: both are server-wide, and ReadTimeout
	// caps the time to read a whole body, which breaks large uploads. ReadHeaderTimeout still cuts
	// off a client that dribbles headers (slowloris), and requestTimeout above bounds handlers.
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	if err = srv.Serve(listener); err != nil {
		return fmt.Errorf("unable to start server. %s", err.Error())
	}

	return nil
}
