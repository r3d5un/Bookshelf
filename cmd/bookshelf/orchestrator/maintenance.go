package orchestrator

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

const RemoveOldScheduledTask string = "Remove Old Scheduled Tasks"

func (m *Module) removeOldScheduledTasks(ctx context.Context) error {
	logger := logging.LoggerFromContext(ctx).With(
		slog.Group(
			"taskRun",
			slog.String("id", uuid.New().String()),
			slog.String("name", RemoveOldScheduledTask),
		),
	)

	logger.Info("removing old completed runs")

	completedState := string(data.CompleteTaskState)
	monthAgo := time.Now().AddDate(0, -1, 0)
	oldScheduledTasks, err := types.ReadAllScheudledTasks(ctx, &m.models, data.Filters{
		Page:      1,
		PageSize:  50_000,
		State:     &completedState,
		RunAtFrom: &monthAgo,
	})
	if err != nil {
		logger.Info("unable to read tasks", "error", err)
		return err
	}

	if oldScheduledTasks.TotalRecords < 1 {
		logger.Info("no old tasks to delete")
		return nil
	}

	for _, scheduledTask := range oldScheduledTasks.Data {
		logger.Info("deleting old task", "scheduledTask", scheduledTask)
		_, err := types.DeleteScheduledTask(ctx, &m.models, scheduledTask.ID)
		if err != nil {
			logger.Error("unable to delete old task", "error", err)
			return err
		}
	}

	return nil
}
