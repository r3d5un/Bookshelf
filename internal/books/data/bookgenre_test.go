package data_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestBookGenreModel(t *testing.T) {
	id := uuid.New()
	title := "TestBookGenreModel"
	description := fmt.Sprintf("This is a test description for %s\n", id.String())
	timestamp := time.Now()
	newBook := data.Book{
		ID:          id,
		Title:       title,
		Description: &description,
		Published:   &timestamp,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	_, err := models.Books.Insert(context.Background(), newBook)
	if err != nil {
		t.Errorf("unable to insert data: %v\n", err)
		return
	}

	id = uuid.New()
	name := "Genre Name TestBookGenreModel"
	description = fmt.Sprintf("Some genre with ID: %s", id.String())
	timestamp = time.Now()
	newGenre := data.Genre{
		ID:          id,
		Name:        &name,
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	_, err = models.Genres.Insert(context.Background(), newGenre)
	if err != nil {
		t.Errorf("unable to insert data: %v\n", err)
		return
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.BookGenres.Insert(context.Background(), newBook.ID, newGenre.ID)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})
}
