package orchestrator

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/Bookshelf/internal/config"
	"github.com/r3d5un/Bookshelf/internal/orchestrator"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/data"
	"github.com/r3d5un/Bookshelf/internal/orchestrator/types"
	"github.com/r3d5un/Bookshelf/internal/system"
)

const ModuleName string = "orchestrator"

type Module struct {
	schedulerID         uuid.UUID
	logger              *slog.Logger
	cfg                 *config.Config
	db                  *pgxpool.Pool
	models              data.Models
	scheduler           *orchestrator.Scheduler
	done                chan struct{}
	taskNotificationCh  chan pgconn.Notification
	taskCollection      orchestrator.Collection
	wg                  sync.WaitGroup
	isSchedulerMasterCh chan bool
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	// TODO: Figure out how to handle coordinate which instance is the scheduler
	// between multiple instances
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

	m.logger.Info("setting scheduler ID from instance")
	m.schedulerID = system.InstanceFromContext(mono.Context())
	m.logger.Info("scheduler instance ID set", "id", m.schedulerID)

	m.logger.Info("injecting database connection")
	m.cfg = mono.Config()

	dbConfig, err := pgxpool.ParseConfig(m.cfg.DB.DSN)
	if err != nil {
		m.logger.Error("unable to parse postgresql pool configuration", "error", err)
		return nil
	}
	m.db, err = pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		m.logger.Error("unable to create connection pool", "error", err)
		return nil
	}
	m.logger.Info("connection pool established")

	m.logger.Info("initializing channels")
	m.wg = sync.WaitGroup{}
	m.done = make(chan struct{})
	m.taskNotificationCh = make(chan pgconn.Notification, 100)

	timeout := time.Duration(m.cfg.DB.Timeout) * time.Second
	m.models = data.NewModels(m.db, &timeout)

	m.logger.Info("creating task collection")
	m.taskCollection = orchestrator.Collection{}
	m.taskCollection.Add("hello-world", m.helloWorld)

	m.logger.Info("creating task runner")
	m.wg.Add(1)
	go m.taskRunner(ctx)

	m.logger.Info("creating task scheduler")
	m.scheduler = orchestrator.NewScheduler(&m.models)
	taskName := "hello-world"
	m.scheduler.AddCronJob(ctx, "* * * * *", types.Task{
		Name: &taskName,
	})
	m.wg.Add(1)
	go func() { // Goroutine for checking the scheduler lock
		defer m.wg.Done()
		m.checkSchedulerLock(ctx)
	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.maintainSchedulerLock(ctx)
	}()

	m.logger.Info("startup complete")

	return nil
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

func (m *Module) Shutdown() {
	m.logger.Info("shutting down module")

	m.logger.Info("stopping scheduler")
	m.scheduler.Stop()
	m.wg.Done()

	m.logger.Info("closing notification channel")
	close(m.taskNotificationCh)

	m.logger.Info("sending stop signal")
	close(m.done)
	close(m.isSchedulerMasterCh)

	m.logger.Info("waiting for scheduler background processes to complete")
	m.wg.Wait()

	// TODO: App hangs upon shutting down the connection pool
	m.logger.Info("closing module connection pool")
	activeConns := m.db.Stat().AcquiredConns()
	if activeConns > 0 {
		m.logger.Error(
			"acquired connections not released, unclean shutdown",
			"activeConnections", activeConns,
		)
	} else {
		m.db.Close()
	}

	m.logger.Info("module shutdown complete")
}

func (m *Module) initModuleLogger(monoLogger *slog.Logger) {
	m.logger = monoLogger.With(slog.Group("module", slog.String("name", ModuleName)))
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
