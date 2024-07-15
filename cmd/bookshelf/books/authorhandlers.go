package books

import (
	"fmt"
	"net/http"

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
