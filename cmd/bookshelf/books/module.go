package books

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "books"

type Module struct {
	logger *slog.Logger
	mux    *http.ServeMux
	db     *sql.DB
	models data.Models
	cfg    *config.Config
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

	m.logger.Info("injecting configuration")
	m.cfg = mono.Config()

	m.logger.Info("injecting mux")
	m.mux = mono.Mux()

	m.logger.Info("injecting database connection")
	m.db = mono.DB()

	m.logger.Info("setting up data models")
	timeout := time.Duration(m.cfg.DB.Timeout) * time.Second
	m.models = data.NewModels(m.db, &timeout)

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
		{"GET /api/v1/books/healthcheck", m.healthcheckHandler},
		// Books
		{"GET /api/v1/books/books", m.ListBookHandler},
		{"GET /api/v1/books/books/{id}", m.GetBookHandler},
		{"POST /api/v1/books/books", m.PostBookHandler},
		{"PATCH /api/v1/books/books/{id}", m.PatchBookHandler},
		{"DELETE /api/v1/books/books/{id}", m.DeleteBookHandler},
		// Authors
		{"GET /api/v1/books/authors", m.ListAuthorHandler},
		{"GET /api/v1/books/authors/{id}", m.GetAuthorHandler},
		{"POST /api/v1/books/authors", m.PostAuthorHandler},
		{"PATCH /api/v1/books/authors/{id}", m.PatchAuthorHandler},
		{"DELETE /api/v1/books/authors/{id}", m.DeleteAuthorHandler},
		// Series
		{"GET /api/v1/books/series", m.ListSeriesHandler},
		{"GET /api/v1/books/series/{id}", m.GetSeriesHandler},
		{"POST /api/v1/books/series", m.PostSeriesHandler},
		{"PATCH /api/v1/books/series/{id}", m.PatchSeriesHandler},
		{"DELETE /api/v1/books/series/{id}", m.DeleteSeriesHandler},
		// Genre
		{"GET /api/v1/books/genre", m.ListGenreHandler},
		{"GET /api/v1/books/genre/{id}", m.GetGenreHandler},
		{"POST /api/v1/books/genre", m.PostGenreHandler},
		{"PATCH /api/v1/books/genre/{id}", m.PatchGenreHandler},
		{"DELETE /api/v1/books/genre/{id}", m.DeleteGenreHandler},
	}

	m.logger.Info("adding protected endpoints")
	for _, d := range scimRouteDefinitions {
		m.logger.Info("adding route", "route", d.Path)
		mux.Handle(d.Path, d.Handler)
	}
}
