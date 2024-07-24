package ui

import (
	"net/http"
	"slices"
	"time"

	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/rest"
)

func (m *Module) Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "home.tmpl", &templateData{})
}

func (m *Module) MyLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "mylibrary.tmpl", &templateData{})
}

func (m *Module) Discover(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "discover.tmpl", &templateData{})
}

func (m *Module) Authors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "authors.tmpl", &templateData{})
}

func (m *Module) AuthorViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "author.tmpl", &templateData{})
}

func (m *Module) CurrentlyReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "currentlyreading.tmpl", &templateData{})
}

func (m *Module) FinishedReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "finishedreading.tmpl", &templateData{})
}

func (m *Module) MyLibraryBookList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	date := time.Now()
	erikson := "Steven Erikson"
	malazan := "Malazan: Book of the Fallen"

	sanderson := "Brandon Sanderson"
	stormlight := "Stormlight Archives"
	cosmere := "Cosmere"

	data := templateData{
		MyLibraryBooks: []myLibraryBook{
			{
				Title:     "Gardens of the Moon",
				Series:    []*string{&malazan},
				Authors:   []*string{&erikson},
				Published: &date,
				Added:     &date,
				Status:    "Read",
			},
			{
				Title:     "Deadhouse Gates",
				Series:    []*string{&malazan},
				Authors:   []*string{&erikson},
				Published: &date,
				Added:     &date,
				Status:    "Want to Read",
			},
			{
				Title:     "The Way of Kings",
				Series:    []*string{&stormlight, &cosmere},
				Authors:   []*string{&sanderson},
				Published: &date,
				Added:     &date,
				Status:    "Dropped",
			},
			{
				Title:     "Words of Radiance",
				Series:    []*string{&stormlight, &cosmere},
				Authors:   []*string{&sanderson},
				Published: &date,
				Added:     &date,
				Status:    "Reading",
			},
		},
	}

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "librarybooklisting.tmpl", &data)
}

func (m *Module) BookViewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "book.tmpl", &templateData{})
}

func (m *Module) DiscoverCategoryMenuHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("reading requested category")
	category, err := rest.ReadStringParam("category", r)
	if err != nil {
		logger.Info("unable to read category parameter", "error", err)
		rest.BadRequestResponse(w, r, "unable to read category parameter")
		return
	}
	logger.Info("category parsed", "category", category)

	data := templateData{SelectedCategory: *category}
	allowedCategories := []string{"books", "authors", "genres"}
	if !slices.Contains(allowedCategories, data.SelectedCategory) {
		logger.Info(
			"requested category not implemented",
			"category", *category,
			"allowedCategories", allowedCategories,
		)
		rest.BadRequestResponse(w, r, "requested category not implemented")
		return
	}

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "discoveryCategoryMenu.tmpl", &data)
}

func (m *Module) DiscoverContentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("reading requested category")
	category, err := rest.ReadStringParam("category", r)
	if err != nil {
		logger.Info("unable to read category parameter", "error", err)
		rest.BadRequestResponse(w, r, "unable to read category parameter")
		return
	}
	logger.Info("category parsed", "category", category)

	logger.Info("rendering UI component", "category", *category)
	switch *category {
	case "books":
		m.renderPartial(w, http.StatusOK, "discoverBooks.tmpl", &templateData{})
	case "genres":
		m.renderPartial(w, http.StatusOK, "discoverGenres.tmpl", &templateData{})
	case "authors":
		m.renderPartial(w, http.StatusOK, "discoverAuthors.tmpl", &templateData{})
	default:
		logger.Info("unable to read category parameter", "error", err)
		rest.BadRequestResponse(w, r, "unable to read category parameter")
		return
	}
}

func (m *Module) BookSeriesAccordionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "bookSeriesAccordion.tmpl", &templateData{})
}
