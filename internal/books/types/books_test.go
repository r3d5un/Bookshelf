package types_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexBookTypes(t *testing.T) {
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

	genreName := "TestGenreName"
	genreRecord := data.Genre{
		ID:          uuid.New(),
		Name:        &genreName,
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}
	_, err := models.Genres.Insert(context.Background(), genreRecord)
	if err != nil {
		t.Errorf("unable to insert genre: %s\n", err)
		return
	}

	seriesName := "TestSeriesName"
	seriesRecord := data.Series{
		ID:          uuid.New(),
		Name:        &seriesName,
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
		ID:          nil,
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

	var insertedBook *data.Book

	t.Run("TestCreateBook", func(t *testing.T) {
		_, err := types.CreateBook(context.Background(), models, book)
		if err != nil {
			t.Errorf("error occurred when registering new book: %s\n", err)
			return
		}
	})

	t.Run("TestReadBook", func(t *testing.T) {
		bookRecord := data.Book{
			ID:          uuid.New(),
			Title:       "TestGetBookTitle",
			Description: &description,
			Published:   &timestamp,
			CreatedAt:   &timestamp,
			UpdatedAt:   &timestamp,
		}
		insertedBook, err = models.Books.Insert(context.Background(), bookRecord)
		if err != nil {
			t.Errorf("unable to insert book: %s\n", err)
			return
		}
		if _, err := types.ReadBook(context.Background(), models, bookRecord.ID); err != nil {
			t.Errorf("error occurred while retrieving book: %s\n", err)
			return
		}
	})

	t.Run("TestReadAllBooks", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
		}
		bookList, err := types.ReadAllBooks(context.Background(), models, filters)
		if err != nil {
			t.Errorf("unable to read books: %s\n", err)
			return
		}
		if len(bookList) < 1 {
			t.Error("no books returned")
			return
		}
	})

	t.Run("TestGetNonExistingBook", func(t *testing.T) {
		if _, err := types.ReadBook(context.Background(), models, uuid.New()); err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				// desired result. caller is tasked with handling the ErrRecordNotFound error
				return
			default:
				t.Errorf("error occurred while retrieving book: %s\n", err)
				return
			}
		}
	})

	t.Run("TestDeleteBook", func(t *testing.T) {
		if err := types.DeleteBook(context.Background(), models, insertedBook.ID); err != nil {
			t.Errorf("unable to delete book: %s\n", err)
			return
		}
	},
	)
}
