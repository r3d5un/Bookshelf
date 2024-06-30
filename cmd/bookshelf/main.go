package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
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

	app := &application{
		logger: logger,
		mux:    http.NewServeMux(),
	}

	app.logger.Info("running module startup")
	app.setupModules(context.Background())

	err := app.serve()
	if err != nil {
		app.logger.Error("unable to start server", "error", err)
		return err
	}

	app.logger.Info("shutting down modules")
	app.shutdownModules()

	app.logger.Info("exiting...")

	return nil
}

func (app *application) setupModules(ctx context.Context) error {
	for _, v := range app.modules {
		if err := v.Startup(ctx, app); err != nil {
			return err
		}
	}

	return nil
}

func (app *application) shutdownModules() error {
	for _, v := range app.modules {
		if err := v.Shutdown(); err != nil {
			return err
		}
	}

	return nil
}
