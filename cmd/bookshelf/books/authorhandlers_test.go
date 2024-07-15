package books_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/r3d5un/Bookshelf/internal/books/types"
)

func TestAuthorHandlers(t *testing.T) {
	authorName := "Brandon Sanderson"
	description := "such books. many words. wow."
	website := "www.brandonsanderson.com"

	newAuthor := types.NewAuthorData{
		Name: "Steven Erikson",
	}
	id, err := types.CreateAuthor(context.Background(), models, newAuthor)
	if err != nil {
		t.Errorf("uanble to insert test data: %s\n", err)
	}

	t.Run("TestPostAuthorHandler", func(t *testing.T) {
		newAuthor := types.NewAuthorData{
			Name:        authorName,
			Description: &description,
			Website:     &website,
		}

		body, err := json.Marshal(newAuthor)
		if err != nil {
			t.Errorf("unable to marhsal author data: %s\n", err)
			return
		}

		postReq := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/bookshelf/authors",
			strings.NewReader(string(body)),
		)
		postReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PostAuthorHandler)
		handler.ServeHTTP(rr, postReq)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf(
				"handler returned wrong error code: got %d, expected %d",
				status,
				http.StatusCreated,
			)
			return
		}
	})

	t.Run("TestGetAuthorHandler", func(t *testing.T) {
		getReq := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/bookshelf/books/authors",
			nil,
		)
		getReq.Header.Set("Content-Type", "application/json")
		getReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.GetAuthorHandler)
		handler.ServeHTTP(rr, getReq)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(
				"handler returned the wrong error code: got %d, expected %d\n",
				status,
				http.StatusOK,
			)
			return
		}
	})

	t.Run("TestListAuthorHandler", func(t *testing.T) {
	})

	t.Run("TestPatchAuthorHandler", func(t *testing.T) {
	})

	t.Run("TestDeleteAuthorHandler", func(t *testing.T) {
	})
}
