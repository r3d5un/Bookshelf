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
	"github.com/r3d5un/Bookshelf/internal/validator"
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

func (m *Module) GetSeriesHandler(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("querying database for obtain series by id")
	series, err := types.ReadSeries(ctx, &m.models, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("series not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get series", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, series, nil)
}

func (m *Module) ListSeriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.ID = rest.ReadQueryUUID(qs, "id", v)
	input.Filters.Name = rest.ReadQueryString(qs, "name", "")
	input.Filters.Description = rest.ReadQueryString(qs, "description", "")
	input.Filters.CreatedAtFrom = rest.ReadQueryDate(qs, "createdAtFrom", v)
	input.Filters.CreatedAtTo = rest.ReadQueryDate(qs, "createdAtTo", v)
	input.Filters.UpdatedAtFrom = rest.ReadQueryDate(qs, "createdAtFrom", v)
	input.Filters.UpdatedAtTo = rest.ReadQueryDate(qs, "createdAtTo", v)

	input.Filters.Page = rest.ReadQueryInt(qs, "page", 1, v)
	input.Filters.PageSize = rest.ReadQueryInt(qs, "page_size", 1_000, v)

	input.Filters.OrderBy = rest.ReadQueryCommaSeperatedString(qs, "order_by", "name")
	input.Filters.OrderBySafeList = []string{
		"id",
		"updated_at",
		"created_at",
		"name",
		"-id",
		"-updated_at",
		"-created_at",
		"-name",
	}
	logger.InfoContext(ctx, "filters set", "filters", input)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		logger.Info("filter validation failed", "validationErrors", v.Errors)
		rest.FailedValidationResponse(w, r, v.Errors)
		return
	}

	logger.Info("getting series", "filters", input.Filters)
	series, err := types.ReadAllSeries(ctx, &m.models, input.Filters)
	if err != nil {
		logger.Error("unable to get series", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, series, nil)
}

func (m *Module) PatchSeriesHandler(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("parsing request body")
	var updateData types.Series
	err = rest.ReadJSON(r, &updateData)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}
	updateData.ID = *id

	updatedSeries, err := types.UpdateSeries(ctx, &m.models, updateData)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("series not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get series", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, updatedSeries, nil)
}
