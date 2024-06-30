package main

import (
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/system"
)

type application struct {
	logger  *slog.Logger
	mux     *http.ServeMux
	modules []system.Module
}

func (app *application) Logger() *slog.Logger {
	return app.logger
}

func (app *application) Mux() *http.ServeMux {
	return app.mux
}
