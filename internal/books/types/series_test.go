package types_test

import (
	"context"
	"testing"

	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexSeriesTypes(t *testing.T) {
	seriesName := "Mistborn"
	description := "Some series description"

	t.Run("TestCreateSeries", func(t *testing.T) {
		newSeries := types.NewSeriesData{
			Name:        seriesName,
			Description: &description,
		}

		_, err := types.CreateSeries(context.Background(), models, newSeries)
		if err != nil {
			t.Errorf("error occurred while registering a new author: %s\n", err)
			return
		}
	})

	t.Run("TestReadSeries", func(t *testing.T) {
	})
}
