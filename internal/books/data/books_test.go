package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestBookModel(t *testing.T) {
	id := uuid.New()
	title := "TestBookTitle"
	description := "This is a test description."
	timestamp := time.Now()
	newBook := data.Book{
		ID:          id,
		Title:       title,
		Description: &description,
		Published:   &timestamp,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.Books.Insert(context.Background(), newBook)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.Books.Get(context.Background(), id)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}
	})

	t.Run("GetAll", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
			ID:       &id,
		}

		_, nRows, err := models.Books.GetAll(context.Background(), filters)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}
		if *nRows < 1 {
			t.Error("no results returned")
			return
		}
	})

	t.Run("Update", func(t *testing.T) {
		newBook.Title = "NewTitle!"

		res, err := models.Books.Update(context.Background(), newBook)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if res.Title != newBook.Title {
			t.Errorf("expected %s, got %s", newBook.Title, res.Title)
			return
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		description := "Upserted description..."
		newBook.Description = &description

		res, err := models.Books.Upsert(context.Background(), newBook)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if *res.Description != *newBook.Description {
			t.Errorf("expected %s, got %s", *newBook.Description, *res.Description)
			return
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.Books.Delete(context.Background(), newBook.ID)
		if err != nil {
			t.Errorf("unable to delete data: %v\n", err)
			return
		}
	})
}
