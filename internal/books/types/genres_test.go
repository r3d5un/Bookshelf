package types_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
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

	t.Run("TestReadAllGenre", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
		}

		genreList, err := types.ReadAllGenre(context.Background(), models, filters)
		if err != nil {
			t.Errorf("unable to read genre: %s\n", err)
			return
		}
		if len(genreList) < 1 {
			t.Errorf("no books returned")
			return
		}
	})

	t.Run("TestUpdateGenre", func(t *testing.T) {
		newDescription := "this text has been updated"
		newGenreData := types.Genre{
			ID:          *id,
			Description: &newDescription,
		}

		_, err := types.UpdateGenre(context.Background(), models, newGenreData)
		if err != nil {
			t.Errorf("unable to update genre: %s\n", err)
			return
		}
	})

	t.Run("TestDeleteGenre", func(t *testing.T) {
		if err := types.DeleteGenre(context.Background(), models, *id); err != nil {
			t.Errorf("unable to delete genre: %s\n", err)
			return
		}
	})
}
