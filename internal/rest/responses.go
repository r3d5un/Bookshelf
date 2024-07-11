package rest

import (
	"fmt"
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/logging"
)

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	message := fmt.Sprintf(
		"the server encountered a problem and could not process your request: %s\n",
		err,
	)
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func LogError(r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Error(
		"an error occurred",
		"request_method", r.Method,
		"request_url", r.URL.String(),
		"error", err,
	)
}

func ErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.InfoContext(ctx, "writing response")
	err := WriteJSON(w, status, ErrorMessage{Message: message}, nil)
	if err != nil {
		logger.ErrorContext(ctx, "error writing response", "error", err)
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	ErrorResponse(w, r, http.StatusBadRequest, message)
}

func Respond(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	content any,
	headers http.Header,
) {
	err := WriteJSON(w, status, content, headers)
	if err != nil {
		logger := logging.LoggerFromContext(r.Context())

		logger.Error("error writing response", "error", err)
		ServerErrorResponse(w, r, err)
		return
	}
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	message := "the requested resource could not be found"
	logger.InfoContext(ctx, "the requested resource could not be found")
	ErrorResponse(w, r, http.StatusNotFound, message)
}
