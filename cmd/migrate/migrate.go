package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/db"
	"github.com/doublehops/dh-go-framework/internal/logga"

	migrate "github.com/doublehops/go-migration"
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

func run() error {
	var args migrate.Action

	flag.StringVar(&args.Action, "action", "", "the intended action")
	flag.StringVar(&args.Name, "name", "", "the name of the migration")
	flag.IntVar(&args.Number, "number", 0, "The number of migrations to run")

	configFile := flag.String("config", "config.json", "Config file to use")
	flag.Parse()

	args = setFlags(args)

	// Setup config.
	cfg, err := config.New(*configFile)
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

	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("there was an error with os.Getwd(). %s", err.Error())
	}

	args.Path = dir + "/migrations"
	args.DB = DB

	err = args.Migrate()
	if err != nil {
		return fmt.Errorf("there was an error initialising migration. %s", err.Error())
	}

	return nil
}

// setFlags will check that the flags received are valid and assign default ones if not supplied.
func setFlags(args migrate.Action) migrate.Action {
	if found := args.IsValidAction(args.Action); !found {
		args.PrintHelp()
	}

	if args.Action == "create" && args.Name == "" {
		args.PrintHelp()
	}

	if args.Action == "up" && args.Number == 0 {
		args.Number = 9999 // run them all if none defined.
	}

	if args.Action == "down" && args.Number == 0 {
		args.Number = 1 // run just one if none defined.
	}

	return args
}
