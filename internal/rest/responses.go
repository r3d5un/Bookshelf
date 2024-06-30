package rest

import (
	"net/http"

	"github.com/r3d5un/Bookshelf/internal/logging"
)

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	logger.InfoContext(ctx, "the server encountered a problem and could not process your request")
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
