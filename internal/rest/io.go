package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/validator"
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

func ReadStringParam(key string, r *http.Request) (*string, error) {
	s := r.PathValue(key)
	if s == "" {
		return nil, fmt.Errorf("empty string parameter")
	}

	return &s, nil
}

func ReadQueryString(
	qs url.Values,
	key string,
	defaultValue string,
) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadQueryStrings(qs url.Values, key string, defaultValues string) string {
	values, ok := qs[key]
	if !ok || len(values) == 0 {
		return defaultValues
	}

	return strings.Join(values, ",")
}

func ReadQueryCommaSeperatedString(
	qs url.Values,
	key string,
	defaultValue string,
) []string {
	s := qs.Get(key)

	if s == "" {
		return []string{defaultValue}
	}

	splitValues := strings.Split(s, ",")

	var seen []string
	var values []string
	for _, val := range splitValues {
		trimmedVal := strings.TrimSpace(val)
		normalizedVal := strings.TrimPrefix(trimmedVal, "-")
		if trimmedVal != "" && !slices.Contains(seen, normalizedVal) {
			seen = append(seen, normalizedVal)
			values = append(values, trimmedVal)
		}
	}

	return values
}

func ReadQueryInt(
	qs url.Values,
	key string,
	defaultValue int,
	v *validator.Validator,
) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func ReadQueryUUID(
	qs url.Values,
	key string,
	v *validator.Validator,
) *uuid.UUID {
	s := qs.Get(key)
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		v.AddError(key, "must be an uuid")
		return nil
	}
	return &id
}

func ReadQueryDate(
	qs url.Values,
	key string,
	v *validator.Validator,
) *time.Time {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, s); err == nil {
			return &date
		}
	}

	v.AddError(key, fmt.Sprintf("not a valid date format, accepting %s", formats))
	return nil
}
