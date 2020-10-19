package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/igomonov88/nimbler_writer/config"
	"github.com/igomonov88/nimbler_writer/internal/platform/database"
	"github.com/igomonov88/nimbler_writer/internal/schema"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Configuration
	cfg, err := config.Parse("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "parsing config")
	}

	// This is used for multiple commands below.
	dbConfig := database.Config{
		User:       cfg.Database.User,
		Password:   cfg.Database.Password,
		Host:       cfg.Database.Host,
		Name:       cfg.Database.Name,
		DisableTLS: cfg.Database.DisableTLS,
	}
	switch os.Args[1] {
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	default:
		err = errors.New("Must specify a command")
	}

	if err != nil {
		return err
	}

	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrations complete")
	return nil
}

func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return err
	}

	fmt.Println("Seed data complete")
	return nil
}
