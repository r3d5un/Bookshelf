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

func TestGenreHandlers(t *testing.T) {
	genreName := "High Fantasy"
	description := "such books. many words. wow."

	newGenre := types.NewGenreData{
		Name: "Steven Erikson",
	}
	id, err := types.CreateGenre(context.Background(), models, newGenre)
	if err != nil {
		t.Errorf("uanble to insert test data: %s\n", err)
	}

	t.Run("TestPostGenreHandler", func(t *testing.T) {
		newGenre := types.NewGenreData{
			Name:        genreName,
			Description: &description,
		}

		body, err := json.Marshal(newGenre)
		if err != nil {
			t.Errorf("unable to marhsal genre data: %s\n", err)
			return
		}

		postReq := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/bookshelf/genre",
			strings.NewReader(string(body)),
		)
		postReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PostGenreHandler)
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

	t.Run("TestGetGenreHandler", func(t *testing.T) {
		getReq := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/bookshelf/books/genre",
			nil,
		)
		getReq.Header.Set("Content-Type", "application/json")
		getReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.GetGenreHandler)
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
}
