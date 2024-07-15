package types_test

import (
	"context"
	"testing"

	// "github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestComplexAuthorTypes(t *testing.T) {
	authorName := "Brandon Sanderson"
	description := "such books. many words. wow."
	website := "www.brandonsanderson.com"
	// newAuthorData := types.Author{
	// 	ID:          uuid.New(),
	// 	Name:        &authorName,
	// 	Description: &description,
	// 	Website:     &website,
	// }

	t.Run("TestCreateAuthor", func(t *testing.T) {
		newAuthor := types.NewAuthorData{
			Name:        authorName,
			Description: &description,
			Website:     &website,
		}

		_, err := types.CreateAuthor(context.Background(), models, newAuthor)
		if err != nil {
			t.Errorf("error occurred while registering a new author: %s\n", err)
			return
		}
	})

	t.Run("TestReadAuthor", func(t *testing.T) {
	})

	t.Run("TestReadAllAuthor", func(t *testing.T) {
	})

	t.Run("TestUpdateAuthor", func(t *testing.T) {
	})

	t.Run("TestDeleteAuthor", func(t *testing.T) {
	})
}
