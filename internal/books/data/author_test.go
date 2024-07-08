package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestAuthorModel(t *testing.T) {
	id := uuid.New()
	name := "John Doe"
	description := "Some author..."
	website := "www.john-doe.xyz"
	timestamp := time.Now()
	newAuthor := data.Author{
		ID:          id,
		Name:        &name,
		Description: &description,
		Website:     &website,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.Authors.Insert(context.Background(), newAuthor)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.Authors.Get(context.Background(), id)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		t.Log("Get successful!")
	})

	t.Run("GetAll", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
			ID:       &id,
		}

		_, nRows, err := models.Authors.GetAll(context.Background(), filters)
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
		updateName := "Jane Doe"
		newAuthor.Name = &updateName

		res, err := models.Authors.Update(context.Background(), newAuthor)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if *res.Name != *newAuthor.Name {
			t.Errorf("expected %s, got %s", *newAuthor.Name, *res.Name)
			return
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		description := "Upserted description..."
		newAuthor.Description = &description

		res, err := models.Authors.Upsert(context.Background(), newAuthor)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if *res.Description != *newAuthor.Description {
			t.Errorf("expected %s, got %s", *newAuthor.Description, *res.Description)
			return
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.Authors.Delete(context.Background(), newAuthor.ID)
		if err != nil {
			t.Errorf("unable to delete data: %v\n", err)
			return
		}
	})
}
