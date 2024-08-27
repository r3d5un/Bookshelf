package orchestrator

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

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
	logger             *slog.Logger
	cfg                *config.Config
	db                 *pgxpool.Pool
	models             data.Models
	scheduler          *orchestrator.Scheduler
	done               chan struct{}
	taskNotificationCh chan pgconn.Notification
	taskCollection     orchestrator.Collection
	wg                 sync.WaitGroup
}

func (m *Module) Startup(ctx context.Context, mono system.Monolith) (err error) {
	// TODO: Figure out how to handle coordinate which instance is the scheduler
	// between multiple instances
	m.initModuleLogger(mono.Logger())
	m.logger.Info("starting module")

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
	m.logger.Info("starting scheduler")
	m.wg.Add(1)
	m.scheduler.Start()

	m.logger.Info("startup complete")
	return nil
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
	go m.models.TaskNotifications.Listen(ctx, m.taskNotificationCh, m.done)

	for task := range m.taskNotificationCh {
		m.logger.Info("received task", "task", task)
		var taskNotification data.TaskNotification
		err := json.Unmarshal([]byte(task.Payload), &taskNotification)
		if err != nil {
			m.logger.Info("unable to decode notification", "error", err)
		}

		// TODO: Consume the task from the queue, not just the notification
		go func() {
			err = m.taskCollection.Run(ctx, taskNotification.Queue)
			if err != nil {
				m.logger.Info("an error occurred while running the task", "error", err)
			}
		}()
	}
}
