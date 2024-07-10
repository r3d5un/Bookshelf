package data_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestBookSeriesModel(t *testing.T) {
	id := uuid.New()
	title := "TestBookSeriesModel"
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
	name := "Series Name TestBookSeriesModel"
	description = fmt.Sprintf("Some series with ID: %s", id.String())
	timestamp = time.Now()
	newSeries := data.Series{
		ID:          id,
		Name:        name,
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	_, err = models.Series.Insert(context.Background(), newSeries)
	if err != nil {
		t.Errorf("unable to insert data: %v\n", err)
		return
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.BookSeries.Insert(context.Background(), newBook.ID, newSeries.ID, 1.0)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})
}
