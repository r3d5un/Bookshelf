package orchestrator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
)

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
				err := m.taskCollection.Run(ctx, notificationPayload.Queue)
				if err != nil {
					m.logger.Info("an error occurred while running the task", "error", err)
				}
			}()

		case <-m.done:
			m.logger.Info("done signal received, stopping task runner")
			return
		}
	}
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

// maintainSchedulerLock is responsible for starting and stopping the scheduler
// based on the state and value of the m.isSchedulerMasterCh.
//
// If the current intance acquires the lock, attempts to maintain the lock will
// occur on each subsequent signal through the m.isSchedulerMasterCh channel.
func (m *Module) maintainSchedulerLock(ctx context.Context) {
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
				m.models.SchedulerLock.MaintainLock(ctx, m.schedulerID)
			}
		}
	}
}
