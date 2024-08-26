package orchestrator

import (
	"context"

	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
	"github.com/robfig/cron/v3"
)

// Scheduler emits tasks to the orchestrator task queue, and notifies
// listeners about new tasks. The scheduler does not run any tasks itself.
type Scheduler struct {
	cron   *cron.Cron
	models *data.Models
}

func NewScheduler(models *data.Models) *Scheduler {
	return &Scheduler{
		cron:   cron.New(),
		models: models,
	}
}

func (s *Scheduler) AddCronJob(ctx context.Context, cronExpr string, task types.Task) (err error) {
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
func (s *Scheduler) Enqueue(ctx context.Context, newTask types.Task) error {
	logger := logging.LoggerFromContext(ctx).With("newTask", newTask)
	newTaskRow := data.TaskQueue{
		Name:     newTask.Name,
		State:    newTask.State,
		RunAt:    newTask.RunAt,
		TaskData: newTask.TaskData,
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

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
