package types_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexGenreTypes(t *testing.T) {
	genreName := "fantasy"
	description := "Here be dragons"

	var id *uuid.UUID

	t.Run("TestCreateGenre", func(t *testing.T) {
		newGenre := types.NewGenreData{
			Name:        genreName,
			Description: &description,
		}

		newGenreID, err := types.CreateGenre(context.Background(), models, newGenre)
		if err != nil {
			t.Errorf("error occurred while registering a new genre: %s\n", err)
			return
		}

		id = newGenreID
	})

	t.Run("TestReadGenre", func(t *testing.T) {
		if _, err := types.ReadGenre(context.Background(), models, *id); err != nil {
			t.Errorf("error occurred while retrieving genre: %s\n", err)
			return
		}
	})
}
