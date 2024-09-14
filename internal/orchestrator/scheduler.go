package orchestrator

import (
	"context"

	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
	"github.com/robfig/cron/v3"
)

// CronScheduler emits tasks to the orchestrator task queue, and notifies
// listeners about new tasks. The scheduler does not run any tasks itself.
type CronScheduler struct {
	cron   *cron.Cron
	models *data.Models
}

func NewScheduler(models *data.Models) *CronScheduler {
	return &CronScheduler{
		cron:   cron.New(),
		models: models,
	}
}

func (s *CronScheduler) AddCronJob(
	ctx context.Context,
	cronExpr string,
	task types.ScheduledTask,
) (err error) {
	logger := logging.LoggerFromContext(ctx)

	_, err = s.cron.AddFunc(cronExpr, func() {
		err := s.Enqueue(ctx, task)
		if err != nil {
			logger.Error("unable to enqueue task", "error", err)
			// Ideally, an alert should be sent here to notify the admin of the error
		}
	})
	if err != nil {
		logger.Error("uanble to add cronjob", "error", err)
		return err
	}

	return nil
}

// TODO: Cleanup the function signature
func (s *CronScheduler) Enqueue(ctx context.Context, newTask types.ScheduledTask) error {
	logger := logging.LoggerFromContext(ctx).With("newTask", newTask)
	newTaskRow := data.TaskQueue{
		Name:  newTask.Name,
		State: newTask.State,
		RunAt: newTask.RunAt,
	}

	enqueuedTask, err := s.models.TaskQueues.Insert(ctx, newTaskRow)
	if err != nil {
		logger.Error("unable to enqueue task")
		return err
	}
	logger.Info("task enqueued", "enqueuedTask", enqueuedTask)

	logger.Info("notifying listeners of new task")
	err = s.models.TaskNotifications.Notify(
		ctx,
		data.TaskNotification{ID: enqueuedTask.ID, Queue: *enqueuedTask.Name},
	)
	logger.Info("listeners notified")

	return nil
}

// Start the scheduler in it's own goroutine, or no-op if already running.
func (s *CronScheduler) Start() {
	s.cron.Start()
}

// Stop stops the scheduler if it is running; otherwise it does nothing.
func (s *CronScheduler) Stop() {
	s.cron.Stop()
}
