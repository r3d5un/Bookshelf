package main

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/system"
)

type application struct {
	logger  *slog.Logger
	mux     *http.ServeMux
	modules []system.Module
	db      *sql.DB
}

func (app *application) Logger() *slog.Logger {
	return app.logger
}

func (app *application) Mux() *http.ServeMux {
	return app.mux
}
