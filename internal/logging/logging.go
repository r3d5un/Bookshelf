package logging

import (
	"context"
	"log/slog"
)

type ContextKey string

const LoggerKey ContextKey = "logger"

// Embeds a logger in the given context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// LoggerFromContext attempts to extract an embedded logger from the
// given context. If no context is found, it returns the default logger
// registered for the application.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
