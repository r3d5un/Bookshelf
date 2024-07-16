package books

import (
	"fmt"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/books/types"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/rest"
)

func (m *Module) PostSeriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing request body")
	var newSeries types.NewSeriesData
	err := rest.ReadJSON(r, &newSeries)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	logger.Info("creating new series", "series", newSeries)
	seriesID, err := types.CreateSeries(ctx, &m.models, newSeries)
	if err != nil {
		logger.Error("unable to create new series records", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("series created", "id", seriesID)

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusCreated, nil, nil)
}
