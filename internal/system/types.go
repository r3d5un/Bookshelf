package system

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/config"
)

type Monolith interface {
	Logger() *slog.Logger
	Mux() *http.ServeMux
	DB() *sql.DB
	Config() *config.Config
	Modules() []Module
}

type Module interface {
	Startup(context.Context, Monolith) error
	Shutdown()
}
