package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestSeriesModel(t *testing.T) {
	id := uuid.New()
	name := "New Series"
	description := "This series is great!"
	timestamp := time.Now()
	newSeries := data.Series{
		ID:          id,
		Name:        name,
		Description: &description,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.Series.Insert(context.Background(), newSeries)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})

	t.Run("Get", func(t *testing.T) {
		_, err := models.Series.Get(context.Background(), id)
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

		_, nRows, err := models.Series.GetAll(context.Background(), filters)
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
		newSeries.Name = "NewNameOfSeries!"

		res, err := models.Series.Update(context.Background(), newSeries)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if res.Name != newSeries.Name {
			t.Errorf("expected %s, got %s", newSeries.Name, res.Name)
			return
		}
	})

	t.Run("Upsert", func(t *testing.T) {
		description := "Upserted description..."
		newSeries.Description = &description

		res, err := models.Series.Upsert(context.Background(), newSeries)
		if err != nil {
			t.Errorf("unable to retrieve result: %v\n", err)
			return
		}

		if *res.Description != *newSeries.Description {
			t.Errorf("expected %s, got %s", *newSeries.Description, *res.Description)
			return
		}
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := models.Series.Delete(context.Background(), newSeries.ID)
		if err != nil {
			t.Errorf("unable to delete data: %v\n", err)
			return
		}
	})
}
