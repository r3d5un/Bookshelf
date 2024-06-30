package main

import "log/slog"

type application struct {
	logger *slog.Logger
}

func (app *application) Logger() *slog.Logger {
	return app.logger
}
