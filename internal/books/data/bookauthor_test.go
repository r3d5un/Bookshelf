package data_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/books/data"
)

func TestBookAuthorModel(t *testing.T) {
	id := uuid.New()
	title := "TestBookAuthorModel"
	description := fmt.Sprintf("This is a test description for %s\n", id.String())
	timestamp := time.Now()
	newBook := data.Book{
		ID:          id,
		Title:       title,
		Description: &description,
		Published:   &timestamp,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	_, err := models.Books.Insert(context.Background(), newBook)
	if err != nil {
		t.Errorf("unable to insert data: %v\n", err)
		return
	}

	id = uuid.New()
	name := "Author Name TestBookAuthorModel"
	description = fmt.Sprintf("Some author with ID: %s", id.String())
	website := "www.john-doe.xyz"
	timestamp = time.Now()
	newAuthor := data.Author{
		ID:          id,
		Name:        &name,
		Description: &description,
		Website:     &website,
		CreatedAt:   &timestamp,
		UpdatedAt:   &timestamp,
	}

	_, err = models.Authors.Insert(context.Background(), newAuthor)
	if err != nil {
		t.Errorf("unable to insert data: %v\n", err)
		return
	}

	t.Run("Insert", func(t *testing.T) {
		_, err := models.BookAuthors.Insert(context.Background(), newBook.ID, newAuthor.ID)
		if err != nil {
			t.Errorf("unable to insert data: %v\n", err)
			return
		}
	})
}
