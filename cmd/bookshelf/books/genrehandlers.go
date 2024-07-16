package books

import (
	"fmt"
	"net/http"

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
