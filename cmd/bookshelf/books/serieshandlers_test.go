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

func TestSeriesHandlers(t *testing.T) {
	seriesName := "High Fantasy"
	description := "such magic. many spells. wow."

	newSeries := types.NewSeriesData{
		Name: "Steven Erikson",
	}
	id, err := types.CreateSeries(context.Background(), models, newSeries)
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

	t.Run("TestGetSeriesHandler", func(t *testing.T) {
		getReq := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/bookshelf/books/series",
			nil,
		)
		getReq.Header.Set("Content-Type", "application/json")
		getReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.GetSeriesHandler)
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

	t.Run("TestListSeriesHandler", func(t *testing.T) {
		baseURL := "/api/v1/bookshelf/books/series"
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

		handler := http.HandlerFunc(mod.ListSeriesHandler)
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

	t.Run("TestPatchSeriesHandler", func(t *testing.T) {
		newDescription := "ThisDescriptionHasBeenUpdated"
		updateData := types.Series{
			Description: &newDescription,
		}
		reqBody, err := json.Marshal(updateData)
		if err != nil {
			t.Errorf("unable to marshal data: %v\n", updateData)
			return
		}

		patchReq := httptest.NewRequest(
			http.MethodPatch,
			"/api/v1/bookshelf/books/series",
			strings.NewReader(string(reqBody)),
		)
		patchReq.Header.Set("Content-Type", "application/json")
		patchReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.PatchSeriesHandler)
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

	t.Run("TestDeleteSeriesHandler", func(t *testing.T) {
		deleteReq := httptest.NewRequest(
			http.MethodDelete,
			"/api/v1/bookshelf/books/series",
			nil,
		)
		deleteReq.Header.Set("Content-Type", "application/json")
		deleteReq.SetPathValue("id", id.String())

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(mod.DeleteSeriesHandler)
		handler.ServeHTTP(rr, deleteReq)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf(
				"handler returned wrong error code: got %d, expected %d",
				status,
				http.StatusNoContent,
			)
			return
		}
	})
}
