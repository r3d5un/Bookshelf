package ui

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/rest"
)

func (m *Module) BookViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	bookID, err := rest.ReadUUIDParam("id", r)
	if err != nil {
		logger.Info("unable to parse parameter", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to parse parameter: %s", err.Error()))
		return
	}
	logger.Info("parameter parsed", "parameter", bookID)

	logger.Info("retrieving book data", "bookId", bookID)
	book, err := m.bookModule.ReadBook(ctx, *bookID)
	if err != nil {
		logger.Error("unable to retrieve data", "error", err, "bookId", bookID)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	logger.Info("setting template data")
	bookData := templateData{
		BookData: *book,
	}
	logger.Info("template data set", "data", bookData)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "book.tmpl", &bookData)
}

func (m *Module) AddAuthorToBookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing book ID from path")
	bookID, err := rest.ReadStringParam("bookID", r)
	if err != nil {
		logger.Info("unable to read category parameter", "error", err)
		rest.BadRequestResponse(w, r, "unable to read category parameter")
		return
	}
	logger.Info("bookID parsed", "bookID", bookID)

	logger.Info("parsing form")
	err = r.ParseForm()
	if err != nil {
		logger.Error("unable to parse form", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	authorName := r.FormValue("modalAuthorNameInput")
	authorID := r.FormValue("modalAuthorIdInput")
	logger.Info("form parsed", "authorName", authorName, "authorID", authorID)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "toast.tmpl", &templateData{})
}

func (m *Module) AddAuthorModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "addAuthorModal.tmpl", &templateData{})
}

func (m *Module) AddAuthorModalDatalist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)
	filters := data.Filters{
		Page:     1,
		PageSize: 50,
	}

	logger.Info("reading authors", "filters", filters)
	authors, err := m.bookModule.ReadAllAuthors(ctx, filters)
	if err != nil {
		logger.Error("error occurred while reading all authors", "filters", filters, "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("authors retrieved", "length", len(authors))

	var buffer bytes.Buffer
	logger.Info("rendering datalist")
	for _, a := range authors {
		buffer.WriteString(
			fmt.Sprintf(`<option value="%s" author-id="%s"></option>`, *a.Name, a.ID.String()),
		)
	}
	logger.Info("datalist rendered")

	logger.Info("responding with UI component")
	m.rawResponse(w, http.StatusOK, buffer.String())
}

func (m *Module) BookSeriesAccordionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	bookID, err := rest.ReadUUIDParam("id", r)
	if err != nil {
		logger.Info("unable to parse parameter", "error", err)
		rest.BadRequestResponse(w, r, fmt.Sprintf("unable to parse parameter: %s", err.Error()))
		return
	}
	logger.Info("parameter parsed", "parameter", bookID)

	book, err := m.bookModule.ReadBook(ctx, *bookID)
	if err != nil {
		logger.Info("uanble to retrieve book data", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	data := templateData{
		SeriesAccordionCollection: []SeriesAccordionCollection{},
	}

	for _, series := range book.Series {
		var seriesAccordion SeriesAccordionCollection

		logger.Info("retrieving books in series", "series", series.ID)
		booksInSeries, err := m.bookModule.ReadBooksBySeries(ctx, series.ID)
		if err != nil {
			logger.Info("unable to retrieve books in series", "error", err)
			rest.ServerErrorResponse(w, r, err)
			return
		}
		logger.Info("retrieved books in series", "books", booksInSeries)

		for orderPlaceholder, book := range booksInSeries {
			logger.Info("constructing series accordion dataset", "series", book.Title)
			seriesaccordion := seriesBookAccordionItem{
				ID:          book.ID.String(),
				Order:       float32(orderPlaceholder),
				Title:       *book.Title,
				Published:   book.Published,
				Description: *book.Description,
				Selected:    false,
			}

			seriesAccordion.Collection = append(seriesAccordion.Collection, seriesaccordion)
		}

		data.SeriesAccordionCollection = append(
			data.SeriesAccordionCollection,
			seriesAccordion,
		)
	}

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "bookSeriesAccordion.tmpl", &data)
}
