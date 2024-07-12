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

	logger.Info("querying database for obtain clients by id")
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
