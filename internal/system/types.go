package system

import (
	"context"
	"log/slog"
	"net/http"
)

type Monolith interface {
	Logger() *slog.Logger
	Mux() *http.ServeMux
}

type Module interface {
	Startup(context.Context, Monolith) error
	Shutdown()
}
