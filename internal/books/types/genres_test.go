package types_test

import (
	"context"
	"testing"

	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexGenreTypes(t *testing.T) {
	genreName := "fantasy"
	description := "Here be dragons"

	t.Run("TestCreateGenre", func(t *testing.T) {
		newGenre := types.NewGenreData{
			Name:        genreName,
			Description: &description,
		}

		_, err := types.CreateGenre(context.Background(), models, newGenre)
		if err != nil {
			t.Errorf("error occurred while registering a new series: %s\n", err)
			return
		}
	})
}
