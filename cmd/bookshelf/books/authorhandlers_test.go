package books_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		baseURL := "/api/v1/bookshelf/books/authors"
		url, err := url.Parse(baseURL)
		if err != nil {
			t.Fatalf("unable to parse URL: %v", err)
			return
		}
		q := url.Query()
		q.Add("page", "1")
		q.Add("pageSize", "10")
		url.RawQuery = q.Encode()

		listReq := httptest.NewRequest(
			http.MethodGet,
			url.String(),
			nil,
		)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.ListAuthorHandler)
		handler.ServeHTTP(rr, listReq)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(
				"handler returned wrong error code: got %d, expected %d",
				status,
				http.StatusOK,
			)
			return
		}
	})

	t.Run("TestPatchAuthorHandler", func(t *testing.T) {
		newDescription := "ThisDescriptionHasBeenUpdated"
		updateData := types.Author{
			Description: &newDescription,
		}
		reqBody, err := json.Marshal(updateData)
		if err != nil {
			t.Errorf("unable to marshal data: %v\n", updateData)
			return
		}

		patchReq := httptest.NewRequest(
			http.MethodPatch,
			"/api/v1/bookshelf/books/authors",
			strings.NewReader(string(reqBody)),
		)
		patchReq.Header.Set("Content-Type", "application/json")
		patchReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PatchAuthorHandler)
		handler.ServeHTTP(rr, patchReq)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf(
				"handler returned wrong error code: got %d, expected %d",
				status,
				http.StatusOK,
			)
			return
		}
	})

	t.Run("TestDeleteAuthorHandler", func(t *testing.T) {
	})
}
