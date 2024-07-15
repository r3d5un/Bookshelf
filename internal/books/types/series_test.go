package types_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexSeriesTypes(t *testing.T) {
	seriesName := "Mistborn"
	description := "Some series description"

	var id *uuid.UUID

	t.Run("TestCreateSeries", func(t *testing.T) {
		newSeries := types.NewSeriesData{
			Name:        seriesName,
			Description: &description,
		}

		newSeriesID, err := types.CreateSeries(context.Background(), models, newSeries)
		if err != nil {
			t.Errorf("error occurred while registering a new series: %s\n", err)
			return
		}

		id = newSeriesID
	})

	t.Run("TestReadSeries", func(t *testing.T) {
		if _, err := types.ReadSeries(context.Background(), models, *id); err != nil {
			t.Errorf("error occurred while retrieving series: %s\n", err)
			return
		}
	})

	t.Run("TestReadAllSeries", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
		}

		seriesList, err := types.ReadAllSeries(context.Background(), models, filters)
		if err != nil {
			t.Errorf("unable to read series: %s\n", err)
			return
		}
		if len(seriesList) < 1 {
			t.Errorf("no books returned")
			return
		}
	})
}
