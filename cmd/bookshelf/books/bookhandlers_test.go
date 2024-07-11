package books_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestBookHandlers(t *testing.T) {
	description := "This is a test description."
	timestamp := time.Now()

	bookRecord := data.Book{
		ID:          uuid.New(),
		Title:       "TestBookTitle",
		Description: &description,
		Published:   &timestamp,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	genreRecord := data.Genre{
		ID:          uuid.New(),
		Name:        "TestGenreName",
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err := models.Genres.Insert(context.Background(), genreRecord)
	if err != nil {
		t.Errorf("unable to insert genre: %s\n", err)
		return
	}

	seriesRecord := data.Series{
		ID:          uuid.New(),
		Name:        "TestSeriesName",
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err = models.Series.Insert(context.Background(), seriesRecord)
	if err != nil {
		t.Errorf("unable to insert series: %s\n", err)
		return
	}

	authorName := "TestAuthorName"
	website := "www.testwebsite.com"
	authorRecord := data.Author{
		ID:          uuid.New(),
		Name:        &authorName,
		Description: &description,
		Website:     &website,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err = models.Authors.Insert(context.Background(), authorRecord)
	if err != nil {
		t.Errorf("unable to insert series: %s\n", err)
		return
	}

	book := types.Book{
		ID:          &bookRecord.ID,
		Title:       &bookRecord.Title,
		Description: bookRecord.Description,
		Published:   bookRecord.Published,
		CreatedAt:   bookRecord.CreatedAt,
		UpdatedAt:   bookRecord.UpdatedAt,
		Authors:     []*data.Author{&authorRecord},
		Genres:      []*data.Genre{&genreRecord},
		BookSeries: []*data.BookSeries{
			{
				BookID:      bookRecord.ID,
				SeriesID:    seriesRecord.ID,
				SeriesOrder: 1.0,
			},
		},
	}

	id, err := types.NewBook(context.Background(), models, book)
	if err != nil {
		t.Errorf("error occurred when registering new book: %s\n", err)
		return
	}

	t.Run("TestPostBookHandler", func(t *testing.T) {
		newBook, err := json.Marshal(book)
		if err != nil {
			t.Errorf("unable to marshal book: %s\n", newBook)
			return
		}

		postReq := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/bookshelf/books",
			nil,
		)
		postReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PostBookHandler)
		handler.ServeHTTP(rr, postReq)
	})

	t.Run("TestGetBook", func(t *testing.T) {
		getReq := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/bookshelf/books",
			nil,
		)
		getReq.Header.Set("Content-Type", "application/json")
		getReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.GetBookHandler)
		handler.ServeHTTP(rr, getReq)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(
				"handler returned the wrong error code: got %d, expected %d\n",
				status,
				http.StatusOK,
			)
			return
		}
	})
}
