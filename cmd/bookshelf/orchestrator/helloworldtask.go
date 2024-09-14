package orchestrator

import (
	"context"
	"errors"

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

	logger, stop := types.NewTaskLogger(ctx, &m.models, TaskName, taskQueueID)
	defer stop()

	logger.InfoContext(ctx, "Hello, World!")

	return nil
}
