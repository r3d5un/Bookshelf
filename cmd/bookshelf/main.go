package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Bookshelf/cmd/bookshelf/books"
	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/database"
	"github.com/r3d5un/Bookshelf/internal/system"
)

func main() {
	if err := run(); err != nil {
		slog.Error("an error occurred", "error", err)
	}

	os.Exit(1)
}

func run() (err error) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	// TODO: Add log group with version and instance ID
	slog.SetDefault(logger)

	logger.Info("starting...")

	logger.Info("loading configuration")
	cfg, configErr := config.New()
	if configErr != nil {
		logger.Error("unable to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("opening database connection pool")

	db, dbErr := database.OpenPool(
		cfg.DB.DSN,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
		time.Duration(cfg.DB.Timeout)*time.Second,
	)
	if dbErr != nil {
		logger.Error("error occurred while creating connection pool", "error", err)
		os.Exit(1)
	}

	app := system.NewMonolith(
		logger,
		http.NewServeMux(),
		[]system.Module{&books.Module{}},
		db,
		cfg,
	)

	app.Logger().Info("running module startup")
	err = app.SetupModules(context.Background())
	if err != nil {
		return err
	}

	err = app.Serve()
	if err != nil {
		app.Logger().Error("unable to start server", "error", err)
		return err
	}

	app.Logger().Info("shutting down modules")
	err = app.ShutdownModules()
	if err != nil {
		return err
	}

	app.Logger().Info("exiting...")

	return nil
}
