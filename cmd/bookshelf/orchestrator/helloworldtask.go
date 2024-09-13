package orchestrator

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

const (
	TaskName string = "Hello, World!"
)

func (m *Module) helloWorld(ctx context.Context) error {

	taskQueueID, ok := ctx.Value("taskQueueID").(uuid.UUID)
	if !ok {
		return errors.New("unable to get task queue ID from context")
	}
	logWriter := types.NewTaskLogWriter(ctx, &m.models, taskQueueID, 100)

	go func() {
		logWriter.LogSink(ctx)
	}()

	handler := slog.NewJSONHandler(io.MultiWriter(&logWriter, os.Stdout), nil)
	logger := slog.New(handler).With(slog.Group(
		"task",
		slog.String("taskName", TaskName),
		slog.String("taskId", taskQueueID.String())),
	)

	logger.InfoContext(ctx, "Hello, World!")

	logWriter.Stop()

	return nil
}
