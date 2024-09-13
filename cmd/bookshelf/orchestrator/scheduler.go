package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/Bookshelf/internal/logging"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
)

// addTasks is where any tasks the scheduler is resposible for queueing
// should be added.
func (m *Module) addTasks(ctx context.Context) error {
	logger := logging.LoggerFromContext(ctx)

	logger.Info("adding tasks")
	tasks := []types.Task{
		types.NewTask("Hello, World!", "* * * * *", false, time.Now(), m.helloWorld),
		types.NewTask(
			RemoveOldScheduledTask, "* * * * *", false, time.Now(), m.removeOldScheduledTasks,
		),
	}

	logger.Info("syncing task with database")
	err := types.SyncTasks(ctx, &m.models, tasks)
	if err != nil {
		logger.Error("an error occurred while syncing tasks with the database", "error", err)
		return err
	}
	logger.Info("tasks synced")

	for _, task := range tasks {
		m.logger.Info("adding task to runner", "task", task)
		m.taskCollection.Add(task.Name, task.Job)

		m.logger.Info("adding task to scheduler", "task", task)
		m.scheduler.AddCronJob(ctx, *task.CronExpr, types.ScheduledTask{
			Name: &task.Name,
		})
	}

	return nil
}

func (m *Module) taskRunner(ctx context.Context) {
	defer m.wg.Done()
	go m.models.TaskNotifications.Listen(ctx, m.taskNotificationCh, m.done)

	for {
		select {
		case notification, ok := <-m.taskNotificationCh:
			if !ok {
				m.logger.Info("task notification channel closed, stopping task runner")
				return
			}

			m.logger.Info("received task", "notification", notification)
			var notificationPayload data.TaskNotification
			if err := json.Unmarshal([]byte(notification.Payload), &notificationPayload); err != nil {
				m.logger.Error("unable to decode notification payload", "error", err)
				continue
			}

			go func() {
				m.runTaskByID(ctx, notificationPayload.ID)
			}()

		case <-m.done:
			m.logger.Info("done signal received, stopping task runner")
			return
		}
	}
}

// runTaskByID runs the given task by it's ID.
//
// The task is locked, then set from a waiting state to a running state, until the task completes
// or fails.
func (m *Module) runTaskByID(ctx context.Context, id uuid.UUID) {
	logger := logging.LoggerFromContext(ctx).With(slog.String("taskId", id.String()))

	logger.Info("claiming task from queue")
	scheduledTask, err := types.ClaimScheduledTaskByID(ctx, &m.models, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			logger.Info(
				"not able to find task; assuming taking by other worker",
				"error", err,
			)
		default:
			logger.Info("error occurred while consuming ID", "error", err)
		}
		return
	}
	logger.Info("scheduled task claimed", "scheduledTask", scheduledTask)

	skippedState := string(data.SkippedTaskState)
	errorState := string(data.ErrorTaskState)
	completeState := string(data.CompleteTaskState)

	logger.Info("checking if task is enabled in the overview")
	task, err := types.ReadTask(ctx, &m.models, *scheduledTask.Name)
	if err != nil {
		logger.Error("unable to read task from overview", "error", err)
		scheduledTask.State = &errorState
		_, err := types.UpdateScheduledTask(ctx, &m.models, *scheduledTask)
		if err != nil {
			logger.Info("unable to set the scheduled task state", "error", err)
		}
		return
	}
	if !*task.Enabled {
		logger.Info("task not enabled; skipping run")
		scheduledTask.State = &skippedState
		_, err := types.UpdateScheduledTask(ctx, &m.models, *scheduledTask)
		if err != nil {
			logger.Info("unable to set the scheduled task state", "error", err)
		}
		return
	}
	logger.Info("task enabled", "task", task)

	logger.Info("embedding task ID in context")
	ctx = context.WithValue(ctx, "taskQueueID", id)

	logger.Info("running scheduled task", "scheduledTask", scheduledTask, "task", task)
	err = m.taskCollection.Run(ctx, *scheduledTask.Name)
	if err != nil {
		m.logger.Info("an error occurred while running the task", "error", err)
		scheduledTask.State = &errorState
		_, err := types.UpdateScheduledTask(ctx, &m.models, *scheduledTask)
		if err != nil {
			logger.Info("unable to set the scheduled task state", "error", err)
		}
		return
	}
	scheduledTask.State = &completeState
	scheduledTask, err = types.UpdateScheduledTask(ctx, &m.models, *scheduledTask)
	if err != nil {
		logger.Info("unable to set the scheduled task state", "error", err)
		return
	}

	logger.Info("scheduled task completed", "scheduledTask", scheduledTask, "task", task)
}

// checkSchedulerLock attempts to acquire the scheduler lock in a continuous loop.
// The current state of the lock is communicated through the m.isSchedulerMasterCh,
// which is responsible for managing the task scheduler
func (m *Module) checkSchedulerLock(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-m.done:
			m.logger.Info("received done signal; stopping scheduler")
			return
		default:
			acquired, err := m.models.SchedulerLock.AcquireLock(ctx, m.schedulerID)
			if err != nil {
				m.logger.Error("error occurred while acquiring scheduler lock", "error", err)
			}
			m.isSchedulerMasterCh <- acquired
		}
	}
}

// manageScheduler is responsible for starting and stopping the scheduler
// based on the state and value of the m.isSchedulerMasterCh.
//
// If the current intance acquires the lock, attempts to maintain the lock will
// occur on each subsequent signal through the m.isSchedulerMasterCh channel.
func (m *Module) manageScheduler(ctx context.Context) {
	for {
		select {
		case <-m.done:
			m.logger.Info("received done signal; no longer maintaining scheduler lock")
			return
		case active, ok := <-m.isSchedulerMasterCh:
			if !ok {
				m.logger.Info("scheduler lock channel closed")
				return
			}
			if !active {
				m.logger.Info("unable to acquire scheduler lock")
				m.scheduler.Stop()
			} else {
				m.logger.Info("scheduler lock acquired")
				m.scheduler.Start()
				err := m.models.SchedulerLock.MaintainLock(ctx, m.schedulerID)
				if err != nil {
					m.logger.Info("unable to maintain scheduler lock", "error", err)
					m.scheduler.Stop()
					continue
				}
			}
		}
	}
}

func (m *Module) taskReminder(ctx context.Context) {
	waitingStateFilter := string(data.WaitingTaskState)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-m.done:
			m.logger.Info("received done signal; stopping scheduler")
			return
		default:
			timestamp := time.Now()
			logger := logging.LoggerFromContext(ctx).
				With(slog.Group(
					"reminderLoop",
					slog.String("id", uuid.New().String()),
					slog.Time("timestamp", timestamp),
				))

			logger.Info("querying for tasks needing reminders")
			staleTasks, err := types.ReadAllScheudledTasks(ctx, &m.models, data.Filters{
				Page:     1,
				PageSize: 500,
				State:    &waitingStateFilter,
				RunAtTo:  &timestamp,
			})
			if err != nil {
				logger.Error("unable to read stale tasks", "error", err)
			}

			for _, scheduledTask := range staleTasks.Data {
				logger.Info(
					"sending reminder notification for scheduled task",
					"scheduledTask", scheduledTask,
				)
				err := m.models.TaskNotifications.Notify(
					ctx,
					data.TaskNotification{ID: scheduledTask.ID, Queue: *scheduledTask.Name},
				)
				logger.Error("unable to send reminder notification", "error", err)
			}
		}
	}
}
