package books

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/rest"
)

func (m *Module) PostGenreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing request body")
	var newGenre types.NewGenreData
	err := rest.ReadJSON(r, &newGenre)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	logger.Info("creating new genre", "genre", newGenre)
	genreID, err := types.CreateGenre(ctx, &m.models, newGenre)
	if err != nil {
		logger.Error("unable to create new genre records", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("genre created", "id", genreID)

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusCreated, nil, nil)
}

func (m *Module) GetGenreHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing ID")
	id, err := rest.ReadUUIDParam("id", r)
	if err != nil {
		logger.Info("unable to read id", "id", id, "error", err)
		rest.NotFoundResponse(w, r)
		return
	}
	logger.Info("ID parsed", slog.String("id", id.String()))

	logger.Info("querying database for obtain genre by id")
	genre, err := types.ReadGenre(ctx, &m.models, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("genre not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get genre", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, genre, nil)
}
