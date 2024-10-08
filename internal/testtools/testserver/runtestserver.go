package testserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/routes"
	"github.com/doublehops/dh-go-framework/internal/service"
)

func RunTestServer() error {
	ctx := context.Background()

	// Setup config.
	cfg, err := config.New("./config_test.json")
	if err != nil {
		return fmt.Errorf("error starting main. %s", err.Error())
	}

	// Setup logger.
	l, err := logga.New(&cfg.Logging)
	if err != nil {
		return fmt.Errorf("error configuring logger. %s", err.Error())
	}

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

	l.Info(ctx, "Starting server on port :8088", nil)

	// todo - This really needs to be replaced with something that allows timeouts.
	err = http.ListenAndServe(":8088", mux) // nolint:gosec // @todo - remove this exception.
	if err != nil {
		return fmt.Errorf("unable to start server. %s", err.Error())
	}

	return nil
}
