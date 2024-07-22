package ui

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "ui"

type Module struct {
	logger        *slog.Logger
	mux           *http.ServeMux
	cfg           *config.Config
	templateCache map[string]*template.Template
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

	m.logger.Info("injecting configuration")
	m.cfg = mono.Config()

	m.logger.Info("loading templates")
	m.templateCache, err = m.newTemplateCache()
	if err != nil {
		m.logger.Error("unable to load tempaltes", "error", err)
		return err
	}

	m.logger.Info("injecting mux")
	m.mux = mono.Mux()
	m.logger.Info("registering routes")
	m.registerEndpoints(m.mux)

	return nil
}

func (m *Module) Shutdown() {
	m.logger.Info("shutting down module", slog.String("module", ModuleName))
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", ModuleName)))
}

type RouteDefinition struct {
	Path    string
	Handler http.HandlerFunc
}

type RouteDefinitionList []RouteDefinition

func (m *Module) registerEndpoints(mux *http.ServeMux) {
	scimRouteDefinitions := RouteDefinitionList{
		{"GET /api/v1/ui/healthcheck", m.healthcheckHandler},
		{"GET /", m.Home},
		{"GET /library", m.MyLibrary},
		{"GET /discover", m.Discover},
		{"GET /authors", m.Authors},
		{"GET /authors/{id}", m.AuthorViewHandler},
		// UI Components
		{"GET /ui/currentlyreading", m.CurrentlyReading},
		{"GET /ui/finishedreading", m.FinishedReading},
		{"POST /ui/librarybooklist", m.MyLibraryBookList},
	}

	m.logger.Info("adding protected endpoints")
	for _, d := range scimRouteDefinitions {
		m.logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}
