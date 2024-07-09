package system

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
)

type Monolith interface {
	Logger() *slog.Logger
	Mux() *http.ServeMux
	DB() *sql.DB
}

type Module interface {
	Startup(context.Context, Monolith) error
	Shutdown()
}
