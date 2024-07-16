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
	_, err := types.CreateGenre(context.Background(), models, newGenre)
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
}
