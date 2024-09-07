package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
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
		types.NewTask("Hello, World!", "* * * * *", false, time.Now()),
	}

	logger.Info("syncing task with database")
	err := types.SyncTasks(ctx, &m.models, tasks)
	if err != nil {
		logger.Error("an error occurred while syncing tasks with the database", "error", err)
		return err
	}
	logger.Info("tasks synced")

	for _, task := range tasks {
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

// runTaskByID consumes the a task from the task queue that corresponds to the given ID,
// runs the task, then dequeues the task.
//
// The queue is marked for updates, and a transaction active througout the task run, and
// will not be invisible for any other worker instance attempting to consume the same task.
func (m *Module) runTaskByID(ctx context.Context, id uuid.UUID) {
	logger := m.logger.With(slog.String("taskId", id.String()))

	var wg sync.WaitGroup
	taskCh := make(chan types.ScheduledTask, 1)
	taskRunResultCh := make(chan error, 1)
	defer close(taskCh)
	defer close(taskRunResultCh)

	wg.Add(1)
	go func() {
		defer wg.Done()

		task, err := types.ClaimScheduledTaskByID(ctx, &m.models, id)
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

		logger.Info("task sent to runner")
		taskCh <- *task

		logger.Info("waiting for runner return signal")
		taskRunErr := <-taskRunResultCh
		logger.Info("result received")

		if taskRunErr != nil {
			logger.Error("task unsuccssful, setting task state to error", slog.Any("error", err))
			task, err = types.SetScheduledTaskState(ctx, &m.models, id, data.ErrorTaskState)
			if err != nil {
				logger.Error(
					"an error occurred while marking the task as failed",
					"error", err,
				)
			}
			return
		}

		task, err = types.SetScheduledTaskState(ctx, &m.models, id, data.CompleteTaskState)
		if err != nil {
			logger.Error(
				"an error occurred while marking the task as complete",
				"error", err,
			)
		}
		return
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-m.done:
			m.logger.Info("done signal received, stopping task runner")
			return
		case scheduledTask, ok := <-taskCh:
			if !ok {
				m.logger.Info("unable run task, negative taskCh signal", "ok", ok)
				taskRunResultCh <- errors.New("unable to run task; negative taskCh signal")
				return
			}

			// Check if the task is enabled. If not, this run is skipped
			task, err := types.ReadTask(ctx, &m.models, *scheduledTask.Name)
			if !*task.Enabled {
				m.logger.Info(
					"task not enabled; skipping run",
					"task", task,
					"scheduledRun", scheduledTask,
				)
				state := string(data.SkippedTaskState)
				scheduledTask.State = &state
				_, err := types.UpdateScheduledTask(ctx, &m.models, scheduledTask)
				taskRunResultCh <- err
				return
			}

			err = m.taskCollection.Run(ctx, *scheduledTask.Name)
			m.logger.Info("an error occurred while running the task", "error", err)
			taskRunResultCh <- err
		}
	}()

	wg.Wait()
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
