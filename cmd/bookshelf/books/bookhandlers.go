package books

import (
	"errors"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/rest"
)

func (m *Module) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing ID")
	id, err := rest.ReadUUIDParam("id", r)
	if err != nil {
		logger.Info("unable to read id", "id", id, "error", err)
		rest.NotFoundResponse(w, r)
		return
	}

	logger.InfoContext(ctx, "querying database for obtain clients by id")
	book, err := m.models.Books.Get(ctx, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.InfoContext(ctx, "clients not found", "client", book)
			rest.NotFoundResponse(w, r)
		default:
			logger.ErrorContext(ctx, "unable to get client", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, book, nil)
}
