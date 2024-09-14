package orchestrator

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

const RemoveOldScheduledTask string = "Remove Old Scheduled Tasks"

func (m *Module) removeOldScheduledTasks(ctx context.Context) error {
	taskQueueID, ok := ctx.Value("taskQueueID").(uuid.UUID)
	if !ok {
		return errors.New("unable to get task queue ID from context")
	}

	logger, stopLogger := types.NewTaskLogger(ctx, &m.models, RemoveOldScheduledTask, taskQueueID)
	defer stopLogger()
	ctx = context.WithValue(ctx, logging.LoggerKey, logger)

	logger.Info("starting task schedule maintenance task")

	monthAgo := time.Now().AddDate(0, -1, 0)
	logger.Info("checking for old tasks", "date", monthAgo)
	oldScheduledTasks, err := types.ReadAllScheudledTasks(ctx, &m.models, data.Filters{
		Page:     1,
		PageSize: 50_000,
		RunAtTo:  &monthAgo,
	})
	if err != nil {
		logger.Info("unable to read tasks", "error", err)
		return err
	}
	if oldScheduledTasks.TotalRecords < 1 {
		logger.Info("no old tasks to delete")
		return nil
	}
	logger.Info("old tasks retrieved", "totalTasks", oldScheduledTasks.TotalRecords)

	logger.Info("removing old completed runs")
	for _, scheduledTask := range oldScheduledTasks.Data {
		logger.Info("deleting old task", "scheduledTask", scheduledTask)
		_, err := types.DeleteScheduledTask(ctx, &m.models, scheduledTask.ID)
		if err != nil {
			logger.Error("unable to delete old task", "error", err)
			return err
		}
	}

	logger.Info("task schedule maintenance complete")

	return nil
}
