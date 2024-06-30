package books

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/system"
)

const moduleName string = "books"

type Module struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module", slog.String("module", moduleName))

	m.logger.Info("injecting mux")
	m.mux = mono.Mux()

	m.logger.Info("registering routes")
	m.registerEndpoints(m.mux)

	return nil
}

func (m *Module) Shutdown() {
	return
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", moduleName)))
}

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (m *Module) registerEndpoints(mux *http.ServeMux) {
	scimRouteDefinitions := RouteDefinitionList{
		{"GET /api/v1/books/healthcheck", m.healthcheckHandler},
	}

	m.logger.Info("adding protected endpoints")
	for _, d := range scimRouteDefinitions {
		m.logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}
