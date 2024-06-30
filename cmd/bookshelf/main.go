package main

import (
	"log/slog"
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
	}

	app.logger.Info("exiting...")

	return nil
}
