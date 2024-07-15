package types_test

import (
	"context"
	"testing"

	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
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

	var id *uuid.UUID

	t.Run("TestCreateAuthor", func(t *testing.T) {
		newAuthor := types.NewAuthorData{
			Name:        authorName,
			Description: &description,
			Website:     &website,
		}

		newAuthorID, err := types.CreateAuthor(context.Background(), models, newAuthor)
		if err != nil {
			t.Errorf("error occurred while registering a new author: %s\n", err)
			return
		}

		id = newAuthorID
	})

	t.Run("TestReadAuthor", func(t *testing.T) {
		if _, err := types.ReadAuthor(context.Background(), models, *id); err != nil {
			t.Errorf("error occurred while retrieving author: %s\n", err)
			return
		}
	})

	t.Run("TestReadAllAuthor", func(t *testing.T) {
		filters := data.Filters{
			Page:     1,
			PageSize: 10,
		}

		authorList, err := types.ReadAllAuthors(context.Background(), models, filters)
		if err != nil {
			t.Errorf("unable to read authors: %s\n", err)
			return
		}
		if len(authorList) < 1 {
			t.Errorf("no books returned")
			return
		}
	})

	t.Run("TestUpdateAuthor", func(t *testing.T) {
	})

	t.Run("TestDeleteAuthor", func(t *testing.T) {
	})
}
