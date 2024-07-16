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

func TestSeriesHandlers(t *testing.T) {
	seriesName := "High Fantasy"
	description := "such magic. many spells. wow."

	newSeries := types.NewSeriesData{
		Name: "Steven Erikson",
	}
	_, err := types.CreateSeries(context.Background(), models, newSeries)
	if err != nil {
		t.Errorf("uanble to insert test data: %s\n", err)
	}

	t.Run("TestPostSeriesHandler", func(t *testing.T) {
		newSeries := types.NewSeriesData{
			Name:        seriesName,
			Description: &description,
		}

		body, err := json.Marshal(newSeries)
		if err != nil {
			t.Errorf("unable to marhsal series data: %s\n", err)
			return
		}

		postReq := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/bookshelf/series",
			strings.NewReader(string(body)),
		)
		postReq.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PostSeriesHandler)
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
