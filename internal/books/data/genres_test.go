package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestGenresModel(t *testing.T) {
	id := uuid.New()
	name := "New Genres"
	description := "This series is great!"
	timestamp := time.Now()
	newGenre := data.Genre{
		ID:          id,
		Name:        name,
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.Genres.Insert(context.Background(), newGenre)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.Genres.Get(context.Background(), id)
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

		_, nRows, err := models.Genres.GetAll(context.Background(), filters)
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
		newGenre.Name = "NewNameOfGenres!"

		res, err := models.Genres.Update(context.Background(), newGenre)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if res.Name != newGenre.Name {
			t.Errorf("expected %s, got %s", newGenre.Name, res.Name)
			return
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		description := "Upserted description..."
		newGenre.Description = &description

		res, err := models.Genres.Upsert(context.Background(), newGenre)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if *res.Description != *newGenre.Description {
			t.Errorf("expected %s, got %s", *newGenre.Description, *res.Description)
			return
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.Genres.Delete(context.Background(), newGenre.ID)
		if err != nil {
			t.Errorf("unable to delete data: %v\n", err)
			return
		}
	})
}
