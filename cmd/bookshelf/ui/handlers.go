package ui

import (
	"net/http"
	"time"

	"github.com/r3d5un/Bookshelf/internal/logging"
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
