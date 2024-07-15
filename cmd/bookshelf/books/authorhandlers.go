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

func (m *Module) PostAuthorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing request body")
	var newAuthor types.NewAuthorData
	err := rest.ReadJSON(r, &newAuthor)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	logger.Info("creating new author", "author", newAuthor)
	authorID, err := types.CreateAuthor(ctx, &m.models, newAuthor)
	if err != nil {
		logger.Error("unable to create new author records", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("author created", "id", authorID)

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusCreated, nil, nil)
}

func (m *Module) GetAuthorHandler(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("querying database for obtain author by id")
	author, err := types.ReadAuthor(ctx, &m.models, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("author not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get author", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, author, nil)
}
