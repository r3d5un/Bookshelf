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
