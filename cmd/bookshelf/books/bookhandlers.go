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
	logger.Info("ID parsed", slog.String("id", id.String()))

	logger.Info("querying database for obtain book by id")
	book, err := types.ReadBook(ctx, &m.models, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("book not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get book", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, book, nil)
}

func (m *Module) ListBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	var input struct {
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.ID = rest.ReadQueryUUID(qs, "id", v)
	input.Filters.Title = rest.ReadQueryString(qs, "title", "")
	input.Filters.Description = rest.ReadQueryString(qs, "description", "")
	input.Filters.PublishedFrom = rest.ReadQueryDate(qs, "publishedFrom", v)
	input.Filters.PublishedTo = rest.ReadQueryDate(qs, "publishedTo", v)
	input.Filters.CreatedAtFrom = rest.ReadQueryDate(qs, "createdAtFrom", v)
	input.Filters.CreatedAtTo = rest.ReadQueryDate(qs, "createdAtTo", v)
	input.Filters.UpdatedAtFrom = rest.ReadQueryDate(qs, "createdAtFrom", v)
	input.Filters.UpdatedAtTo = rest.ReadQueryDate(qs, "createdAtTo", v)

	input.Filters.Page = rest.ReadQueryInt(qs, "page", 1, v)
	input.Filters.PageSize = rest.ReadQueryInt(qs, "page_size", 1_000, v)

	input.Filters.OrderBy = rest.ReadQueryCommaSeperatedString(qs, "order_by", "published")
	input.Filters.OrderBySafeList = []string{
		"id",
		"updated_at",
		"created_at",
		"published",
		"title",
		"-id",
		"-updated_at",
		"-created_at",
		"-published",
		"-title",
	}
	logger.InfoContext(ctx, "filters set", "filters", input)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		logger.Info("filter validation failed", "validationErrors", v.Errors)
		rest.FailedValidationResponse(w, r, v.Errors)
		return
	}

	logger.Info("getting books", "filters", input.Filters)
	books, err := types.ReadAllBooks(ctx, &m.models, input.Filters)
	if err != nil {
		logger.Error("unable to get books", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, books, nil)
}

func (m *Module) PostBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing request body")
	var newBook types.Book
	err := rest.ReadJSON(r, &newBook)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	_, err = types.CreateBook(ctx, &m.models, newBook)
	if err != nil {
		logger.Error("unable to create new book records", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusCreated, nil, nil)
}

func (m *Module) PatchBookHandler(w http.ResponseWriter, r *http.Request) {
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
	updateData := types.Book{
		ID: id,
	}
	err = rest.ReadJSON(r, &updateData)
	if err != nil {
		logger.Info("unable to read request body", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to read request body: %s\n", err))
		return
	}

	updatedBook, err := types.UpdateBook(ctx, &m.models, updateData)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("book not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to get book", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusOK, updatedBook, nil)
}

func (m *Module) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
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

	logger.Info("deleting book", "id", id)
	if err := types.DeleteBook(ctx, &m.models, *id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info("book not found", "id", id)
			rest.NotFoundResponse(w, r)
		default:
			logger.Error("unable to delete book", "id", id, "error", err)
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}
	logger.Info("book deleted")

	logger.Info("writing response")
	rest.Respond(w, r, http.StatusNoContent, nil, nil)
}
