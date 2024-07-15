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

func (m *Module) ListAuthorHandler(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("getting authors", "filters", input.Filters)
	authors, err := types.ReadAllAuthors(ctx, &m.models, input.Filters)
	if err != nil {
		logger.Error("unable to get authors", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, authors, nil)
}

func (m *Module) PatchAuthorHandler(w http.ResponseWriter, r *http.Request) {
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
	updateData := types.Author{
		ID: *id,
	}
	err = rest.ReadJSON(r, &updateData)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	updatedAuthor, err := types.UpdateAuthor(ctx, &m.models, updateData)
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
	rest.Respond(w, r, http.StatusOK, updatedAuthor, nil)
}
