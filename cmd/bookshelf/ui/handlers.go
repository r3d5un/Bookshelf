package ui

import (
	"net/http"
	"slices"
	"time"

	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
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

func (m *Module) Series(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "series.tmpl", &templateData{})
}

func (m *Module) NewSeriesModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.renderPartial(w, http.StatusOK, "newSeriesModal.tmpl", &templateData{})
}

func (m *Module) ParseNewSeriesForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing form")
	err := r.ParseForm()
	if err != nil {
		logger.Error("unable to parse form", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	description := r.FormValue("seriesDescriptionTextarea")
	newSeries := types.NewSeriesData{
		Name:        r.FormValue("seriesNameInput"),
		Description: &description,
	}
	logger.Info("form parsed", "newSeries", newSeries)

	logger.Info("creating new series")
	newSeriesID, err := m.bookModule.CreateSeries(ctx, newSeries)
	if err != nil {
		logger.Error("error occurred while creating new series", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("new series created", "id", newSeriesID)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "toast.tmpl", &templateData{})
}

func (m *Module) NewAuthorModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.renderPartial(w, http.StatusOK, "newAuthorModal.tmpl", &templateData{})
}

func (m *Module) ParseNewAuthorForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	err := r.ParseForm()
	if err != nil {
		logger.Error("unable to parse form", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	description := r.FormValue("authorDescriptionTextarea")
	website := r.FormValue("authorWebsiteInput")
	newAuthor := types.NewAuthorData{
		Name:        r.FormValue("authorNameInput"),
		Description: &description,
		Website:     &website,
	}
	logger.Info("form parsed", "newAuthorData", newAuthor)

	logger.Info("creating new author")
	newAuthorID, err := m.bookModule.CreateAuthor(ctx, newAuthor)
	if err != nil {
		logger.Error("error occurred while creating new author", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("new author created", "id", newAuthorID)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "toast.tmpl", &templateData{})
}

func (m *Module) NewGenreModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering page")
	m.render(w, http.StatusOK, "newGenreModal.tmpl", &templateData{})
}

func (m *Module) ParseNewGenreForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing form")
	err := r.ParseForm()
	if err != nil {
		logger.Error("unable to parse form", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	description := r.FormValue("genreDescriptionTextarea")
	newGenre := types.NewGenreData{
		Name:        r.FormValue("genreNameInput"),
		Description: &description,
	}
	logger.Info("form parsed", "newGenre", newGenre)

	logger.Info("creating new author")
	newGenreID, err := m.bookModule.CreateGenre(ctx, newGenre)
	if err != nil {
		logger.Error("error occurred while creating new genre", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("new genre created", "id", newGenreID)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "toast.tmpl", &templateData{})
}

func (m *Module) NewBookModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "newBookModal.tmpl", &templateData{})
}

func (m *Module) ParseNewBookForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("parsing form")
	err := r.ParseForm()
	if err != nil {
		logger.Error("unable to parse form", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}

	description := r.FormValue("bookDescriptionTextarea")
	title := r.FormValue("bookTitleInput")
	timestamp := time.Now()
	newGenre := types.Book{
		Title:       &title,
		Description: &description,
		Published:   &timestamp,
	}
	logger.Info("form parsed", "newBook", newGenre)

	logger.Info("creating new book")
	newBookID, err := m.bookModule.CreateBook(ctx, newGenre)
	if err != nil {
		logger.Error("error occurred while creating new book", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("new book created", "id", newBookID)

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "toast.tmpl", &templateData{})
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

	filters := data.Filters{Page: 1, PageSize: 50}
	logger.Info("retrieving data", "filters", filters)
	books, err := m.bookModule.ReadAllBook(ctx, filters)
	if err != nil {
		logger.Info("error occurred while retrieving books", "error", err)
		rest.ServerErrorResponse(w, r, err)
		return
	}
	logger.Info("data retrieved", "books", books)

	placeholderData := templateData{
		MyLibraryBooks: books,
	}

	logger.Info("rendering UI component")
	m.renderPartial(w, http.StatusOK, "librarybooklisting.tmpl", &placeholderData)
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
