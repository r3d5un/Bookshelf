package ui

import (
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/logging"
)

func (m *Module) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "home.tmpl", &templateData{})
}

func (m *Module) MyLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "mylibrary.tmpl", &templateData{})
}

func (m *Module) Discover(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "discover.tmpl", &templateData{})
}

func (m *Module) TestHTMX(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("button press registered", "request", r)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})

func (m *Module) CurrentlyReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "currentlyreading.tmpl", &templateData{})
}

func (m *Module) FinishedReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "finishedreading.tmpl", &templateData{})
}
