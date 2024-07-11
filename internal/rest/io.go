package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadJSON(r *http.Request, data interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func ReadUUIDParam(key string, r *http.Request) (*uuid.UUID, error) {
	rawID := r.PathValue(key)
	if rawID == "" {
		return nil, fmt.Errorf("%s UUID parameter is emtpy", key)
	}

	id, err := uuid.Parse(rawID)
	if err != nil {
		return nil, fmt.Errorf("%s uuid parameter contains an invalid value: %s\n", key, rawID)
	}

	return &id, nil
}
